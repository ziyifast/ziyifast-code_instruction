package xa

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// XATransactionCoordinator XA事务协调者
type XATransactionCoordinator struct {
	participants []XAParticipant
	mutex        sync.RWMutex
	txID         string
}

// XAParticipant XA事务参与者接口
type XAParticipant interface {
	Start(ctx context.Context) error
	ExecuteBusinessLogic(ctx context.Context) error
	End(ctx context.Context) error
	Prepare(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetID() string
}

// NewXATransactionCoordinator 创建新的XA事务协调者
func NewXATransactionCoordinator(txID string) *XATransactionCoordinator {
	return &XATransactionCoordinator{
		participants: make([]XAParticipant, 0),
		txID:         txID,
	}
}

// AddParticipant 添加XA事务参与者
func (c *XATransactionCoordinator) AddParticipant(participant XAParticipant) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.participants = append(c.participants, participant)
}

// ExecuteTwoPhaseCommit 执行两阶段提交
func (c *XATransactionCoordinator) ExecuteTwoPhaseCommit(ctx context.Context) error {
	// 设置默认超时时间
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 第一阶段：准备阶段
	if err := c.preparePhase(ctxWithTimeout); err != nil {
		// 任一参与者准备失败，全部回滚
		c.rollbackAll(ctxWithTimeout)
		return fmt.Errorf("prepare phase failed: %w", err)
	}

	// 第二阶段：提交阶段
	if err := c.commitPhase(ctxWithTimeout); err != nil {
		// 注意：提交阶段失败需要人工干预
		log.Printf("CRITICAL: commit phase failed, manual intervention required for transaction %s: %v", c.txID, err)
		return fmt.Errorf("commit phase failed: %w", err)
	}

	return nil
}

// preparePhase 准备阶段
func (c *XATransactionCoordinator) preparePhase(ctx context.Context) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	log.Printf("Starting prepare phase for transaction %s with %d participants", c.txID, len(c.participants))

	for i, participant := range c.participants {
		// 执行业务逻辑：XA Start + 业务SQL
		if err := participant.ExecuteBusinessLogic(ctx); err != nil {
			return fmt.Errorf("participant %s business logic execution failed at index %d: %w", participant.GetID(), i, err)
		}

		// 执行XA END
		if err := participant.End(ctx); err != nil {
			return fmt.Errorf("participant %s END failed at index %d: %w", participant.GetID(), i, err)
		}

		// 执行XA PREPARE
		if err := participant.Prepare(ctx); err != nil {
			return fmt.Errorf("participant %s PREPARE failed at index %d: %w", participant.GetID(), i, err)
		}

		log.Printf("Participant %s prepared successfully", participant.GetID())
	}

	log.Printf("Prepare phase completed successfully for transaction %s", c.txID)
	return nil
}

// commitPhase 提交阶段
func (c *XATransactionCoordinator) commitPhase(ctx context.Context) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	log.Printf("Starting commit phase for transaction %s", c.txID)

	for i, participant := range c.participants {
		if err := participant.Commit(ctx); err != nil {
			// 记录已提交的参与者索引，便于追踪
			log.Printf("CRITICAL: Participant %s commit failed at index %d. Already committed participants may require manual recovery.", participant.GetID(), i)
			return fmt.Errorf("participant %s COMMIT failed at index %d: %w", participant.GetID(), i, err)
		}
		log.Printf("Participant %s committed successfully", participant.GetID())
	}

	log.Printf("Commit phase completed successfully for transaction %s", c.txID)
	return nil
}

// rollbackAll 回滚所有参与者
func (c *XATransactionCoordinator) rollbackAll(ctx context.Context) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	log.Printf("Rolling back all participants for transaction %s", c.txID)

	for i, participant := range c.participants {
		if err := participant.Rollback(ctx); err != nil {
			log.Printf("Participant %s rollback failed at index %d: %v", participant.GetID(), i, err)
		} else {
			log.Printf("Participant %s rolled back successfully", participant.GetID())
		}
	}
}

// BusinessOperation 业务操作函数类型
type BusinessOperation func(*sql.Tx) error

// DatabaseParticipant 数据库参与者
type DatabaseParticipant struct {
	db         *sql.DB
	xaTxID     string
	id         string
	operations []BusinessOperation
	currentTx  *sql.Tx
}

// NewDatabaseParticipant 创建新的数据库参与者
func NewDatabaseParticipant(db *sql.DB, xaTxID, id string) *DatabaseParticipant {
	return &DatabaseParticipant{
		db:         db,
		xaTxID:     xaTxID,
		id:         id,
		operations: make([]BusinessOperation, 0),
	}
}

// AddOperation 添加业务操作
func (d *DatabaseParticipant) AddOperation(op BusinessOperation) {
	d.operations = append(d.operations, op)
}

// Start 开始XA事务
func (d *DatabaseParticipant) Start(ctx context.Context) error {
	query := fmt.Sprintf("XA START '%s'", d.xaTxID)
	log.Printf("Executing: %s", query)
	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("XA START failed for participant %s: %v", d.id, err)
		return fmt.Errorf("XA START failed: %w", err)
	}
	log.Printf("XA START succeeded for participant %s", d.id)
	return nil
}

// ExecuteBusinessLogic 执行业务逻辑
func (d *DatabaseParticipant) ExecuteBusinessLogic(ctx context.Context) error {
	// 开始XA事务
	if err := d.Start(ctx); err != nil {
		return err
	}

	// 创建一个普通事务来执行业务SQL
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	d.currentTx = tx

	// 执行所有业务操作
	for i, op := range d.operations {
		if err := op(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("business operation %d failed: %w", i, err)
		}
	}

	// 注意：这里不提交事务，因为XA事务会在后续阶段处理
	log.Printf("Business logic executed successfully for participant %s", d.id)
	return nil
}

// End 结束XA事务
func (d *DatabaseParticipant) End(ctx context.Context) error {
	// 提交内部事务
	if d.currentTx != nil {
		if err := d.currentTx.Commit(); err != nil {
			log.Printf("Internal transaction commit failed for participant %s: %v", d.id, err)
			return fmt.Errorf("internal transaction commit failed: %w", err)
		}
		d.currentTx = nil
	}

	query := fmt.Sprintf("XA END '%s'", d.xaTxID)
	log.Printf("Executing: %s", query)
	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("XA END failed for participant %s: %v", d.id, err)
		return fmt.Errorf("XA END failed: %w", err)
	}
	log.Printf("XA END succeeded for participant %s", d.id)
	return nil
}

// Prepare 准备提交XA事务
func (d *DatabaseParticipant) Prepare(ctx context.Context) error {
	query := fmt.Sprintf("XA PREPARE '%s'", d.xaTxID)
	log.Printf("Executing: %s", query)
	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("XA PREPARE failed for participant %s: %v", d.id, err)
		return fmt.Errorf("XA PREPARE failed: %w", err)
	}
	log.Printf("XA PREPARE succeeded for participant %s", d.id)
	return nil
}

// Commit 提交XA事务
func (d *DatabaseParticipant) Commit(ctx context.Context) error {
	query := fmt.Sprintf("XA COMMIT '%s'", d.xaTxID)
	log.Printf("Executing: %s", query)
	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("XA COMMIT failed for participant %s: %v", d.id, err)
		return fmt.Errorf("XA COMMIT failed: %w", err)
	}
	log.Printf("XA COMMIT succeeded for participant %s", d.id)
	return nil
}

// Rollback 回滚XA事务
func (d *DatabaseParticipant) Rollback(ctx context.Context) error {
	query := fmt.Sprintf("XA ROLLBACK '%s'", d.xaTxID)
	log.Printf("Executing: %s", query)
	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("XA ROLLBACK failed for participant %s: %v", d.id, err)
		return fmt.Errorf("XA ROLLBACK failed: %w", err)
	}

	// 回滚内部事务（如果存在）
	if d.currentTx != nil {
		d.currentTx.Rollback()
		d.currentTx = nil
	}

	log.Printf("XA ROLLBACK succeeded for participant %s", d.id)
	return nil
}

// GetID 获取参与者ID
func (d *DatabaseParticipant) GetID() string {
	return d.id
}

// UsageExample 使用示例
func UsageExample() {
	// 创建XA事务协调者
	coordinator := NewXATransactionCoordinator("tx_12345")
	fmt.Println("Usage example:", coordinator)

	// 创建数据库连接（示例）
	// db1, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db1")
	// db2, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db2")

	// 创建参与者
	// participant1 := NewDatabaseParticipant(db1, "tx_12345", "account_service")
	// participant2 := NewDatabaseParticipant(db2, "tx_12345", "inventory_service")

	// 添加业务操作
	// participant1.AddOperation(func(tx *sql.Tx) error {
	//     _, err := tx.Exec("UPDATE account SET balance = balance - 100 WHERE user_id = 1")
	//     return err
	// })
	//
	// participant2.AddOperation(func(tx *sql.Tx) error {
	//     _, err := tx.Exec("UPDATE inventory SET stock = stock - 1 WHERE product_id = 1")
	//     return err
	// })

	// 添加参与者到协调者
	// coordinator.AddParticipant(participant1)
	// coordinator.AddParticipant(participant2)

	// 执行两阶段提交
	// ctx := context.Background()
	// if err := coordinator.ExecuteTwoPhaseCommit(ctx); err != nil {
	//     log.Fatalf("XA transaction failed: %v", err)
	// }

	// fmt.Println("XA transaction completed successfully")
}
