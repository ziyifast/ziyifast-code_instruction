package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// MessageStatus 消息状态枚举
type MessageStatus int

const (
	MessageStatusPending MessageStatus = iota // 待发送
	MessageStatusSent                         // 已发送
	MessageStatusFailed                       // 发送失败
)

// LocalMessage 本地消息表结构
type LocalMessage struct {
	ID         int64         `json:"id"`
	MessageID  string        `json:"message_id"`  // 消息唯一标识
	Topic      string        `json:"topic"`       // 消息主题
	Content    string        `json:"content"`     // 消息内容
	Status     MessageStatus `json:"status"`      // 消息状态
	RetryCount int           `json:"retry_count"` // 重试次数
	CreatedAt  time.Time     `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time     `json:"updated_at"`  // 更新时间
}

// MessageProducer 消息生产者接口
type MessageProducer interface {
	SendMessage(topic string, content []byte) error
}

// LocalMessageService 本地消息服务
type LocalMessageService struct {
	db        *sql.DB
	producer  MessageProducer
	batchSize int
}

// NewLocalMessageService 创建本地消息服务
func NewLocalMessageService(db *sql.DB, producer MessageProducer) *LocalMessageService {
	service := &LocalMessageService{
		db:        db,
		producer:  producer,
		batchSize: 100,
	}

	// 初始化消息表
	service.initTable()
	return service
}

// initTable 初始化消息表
func (s *LocalMessageService) initTable() {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS local_messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        message_id TEXT UNIQUE NOT NULL,
        topic TEXT NOT NULL,
        content TEXT NOT NULL,
        status INTEGER DEFAULT 0,
        retry_count INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := s.db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create local_messages table:", err)
	}
}

// SaveMessageInTransaction 在业务事务中保存消息
func (s *LocalMessageService) SaveMessageInTransaction(tx *sql.Tx, messageID, topic string, content []byte) error {
	insertSQL := `
    INSERT INTO local_messages (message_id, topic, content, status, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?)`

	_, err := tx.Exec(insertSQL, messageID, topic, string(content), MessageStatusPending,
		time.Now(), time.Now())
	return err
}

// ProcessPendingMessages 处理待发送消息
func (s *LocalMessageService) ProcessPendingMessages(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.sendPendingMessages()
		}
	}
}

// sendPendingMessages 发送待处理消息
func (s *LocalMessageService) sendPendingMessages() {
	// 查询待发送消息
	querySQL := `
    SELECT id, message_id, topic, content, retry_count
    FROM local_messages 
    WHERE status = ? AND retry_count < 3
    ORDER BY created_at ASC
    LIMIT ?`

	rows, err := s.db.Query(querySQL, MessageStatusPending, s.batchSize)
	if err != nil {
		log.Printf("Query pending messages failed: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var msg LocalMessage
		err := rows.Scan(&msg.ID, &msg.MessageID, &msg.Topic, &msg.Content, &msg.RetryCount)
		if err != nil {
			log.Printf("Scan message failed: %v", err)
			continue
		}

		// 发送消息
		if err := s.sendMessage(&msg); err != nil {
			s.handleSendFailure(&msg)
		} else {
			s.handleSendSuccess(&msg)
		}
	}
}

// sendMessage 发送单条消息
func (s *LocalMessageService) sendMessage(message *LocalMessage) error {
	return s.producer.SendMessage(message.Topic, []byte(message.Content))
}

// handleSendSuccess 处理发送成功
func (s *LocalMessageService) handleSendSuccess(message *LocalMessage) {
	updateSQL := `
    UPDATE local_messages 
    SET status = ?, updated_at = ?
    WHERE id = ?`

	_, err := s.db.Exec(updateSQL, MessageStatusSent, time.Now(), message.ID)
	if err != nil {
		log.Printf("Update message status failed: %v", err)
	} else {
		log.Printf("Message sent successfully: %s", message.MessageID)
	}
}

// handleSendFailure 处理发送失败
func (s *LocalMessageService) handleSendFailure(message *LocalMessage) {
	updateSQL := `
    UPDATE local_messages 
    SET retry_count = retry_count + 1, updated_at = ?
    WHERE id = ?`

	_, err := s.db.Exec(updateSQL, time.Now(), message.ID)
	if err != nil {
		log.Printf("Update message retry count failed: %v", err)
	} else {
		log.Printf("Message send failed, retry count: %d, message: %s",
			message.RetryCount+1, message.MessageID)
	}
}

// MockMessageProducer 模拟消息生产者实现
type MockMessageProducer struct{}

func (m *MockMessageProducer) SendMessage(topic string, content []byte) error {
	// 模拟偶发性发送失败
	if time.Now().Unix()%10 == 0 {
		return fmt.Errorf("network error")
	}

	log.Printf("Message sent to topic '%s': %s", topic, string(content))
	return nil
}

// BusinessService 业务服务示例
type BusinessService struct {
	db      *sql.DB
	message *LocalMessageService
}

// NewBusinessService 创建业务服务
func NewBusinessService(db *sql.DB, message *LocalMessageService) *BusinessService {
	return &BusinessService{
		db:      db,
		message: message,
	}
}

// ProcessOrder 处理订单业务（包含本地消息表）
func (b *BusinessService) ProcessOrder(orderID string, userID string, amount float64) error {
	// 开启数据库事务
	tx, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback()

	// 1. 业务处理：创建订单记录
	insertOrderSQL := `INSERT INTO orders (order_id, user_id, amount) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertOrderSQL, orderID, userID, amount)
	if err != nil {
		return fmt.Errorf("create order failed: %w", err)
	}

	// 2. 在同一事务中保存消息
	messageContent := map[string]interface{}{
		"order_id": orderID,
		"user_id":  userID,
		"amount":   amount,
		"event":    "order_created",
	}

	contentBytes, _ := json.Marshal(messageContent)
	messageID := fmt.Sprintf("msg_%s_%d", orderID, time.Now().Unix())

	err = b.message.SaveMessageInTransaction(tx, messageID, "order_events", contentBytes)
	if err != nil {
		return fmt.Errorf("save message failed: %w", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	log.Printf("Order processed successfully: %s", orderID)
	return nil
}

// 初始化数据库和表结构
func initDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./local_message.db")
	if err != nil {
		log.Fatal("Open database failed:", err)
	}

	// 创建订单表示例
	createOrderTableSQL := `
    CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        order_id TEXT UNIQUE NOT NULL,
        user_id TEXT NOT NULL,
        amount REAL NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(createOrderTableSQL)
	if err != nil {
		log.Fatal("Create orders table failed:", err)
	}

	return db
}

// Example 使用示例
func Example() {
	// 初始化数据库
	db := initDatabase()
	defer db.Close()

	// 创建消息生产者和服务
	producer := &MockMessageProducer{}
	messageService := NewLocalMessageService(db, producer)

	// 创建业务服务
	businessService := NewBusinessService(db, messageService)

	// 启动消息处理器
	ctx, cancel := context.WithCancel(context.Background())
	go messageService.ProcessPendingMessages(ctx)
	defer cancel()

	// 处理订单业务
	err := businessService.ProcessOrder("ORDER_001", "USER_001", 99.99)
	if err != nil {
		log.Printf("Process order failed: %v", err)
	}

	// 等待一段时间观察消息发送
	time.Sleep(10 * time.Second)
}

func main() {
	Example()
}
