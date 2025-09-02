package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentStatus 定义支付状态
type PaymentStatus string

const (
	StatusProcessing PaymentStatus = "PROCESSING"
	StatusSuccess    PaymentStatus = "SUCCESS"
	StatusFailed     PaymentStatus = "FAILED"
)

// PaymentRequest 支付请求结构
type PaymentRequest struct {
	IdempotentToken string  `json:"idempotent_token"`
	Amount          float64 `json:"amount"`
	UserID          string  `json:"user_id"`
	Description     string  `json:"description"`
}

// PaymentResult 支付结果结构
type PaymentResult struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	OrderID   string `json:"order_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// PaymentRecord 支付记录结构
type PaymentRecord struct {
	IdempotentKey string         `json:"idempotent_key"`
	Status        PaymentStatus  `json:"status"`
	Result        *PaymentResult `json:"result,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	ErrorMessage  string         `json:"error_message,omitempty"`
}

// RedisLockService Redis分布式锁服务
type RedisLockService struct {
	client redis.Cmdable
	ctx    context.Context
}

func NewRedisLockService(client redis.Cmdable) *RedisLockService {
	return &RedisLockService{
		client: client,
		ctx:    context.Background(),
	}
}

// generateLockValue 生成锁的唯一标识值
func (r *RedisLockService) generateLockValue() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// AcquireLock 获取分布式锁
func (r *RedisLockService) AcquireLock(key string, expiration time.Duration) (string, bool) {
	value, err := r.generateLockValue()
	if err != nil {
		return "", false
	}

	// 使用SET命令的NX和EX选项原子性地获取锁
	success, err := r.client.SetNX(key, value, expiration).Result()
	if err != nil {
		return "", false
	}

	if !success {
		return "", false
	}

	return value, true
}

// TryAcquireLock 尝试获取分布式锁，支持等待时间
func (r *RedisLockService) TryAcquireLock(key string, waitTime, expiration time.Duration) (string, bool) {
	value, acquired := r.AcquireLock(key, expiration)
	if acquired {
		return value, true
	}

	// 如果获取锁失败，则等待并重试
	deadline := time.Now().Add(waitTime)
	for time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond) // 短暂休眠避免过度竞争
		value, acquired := r.AcquireLock(key, expiration)
		if acquired {
			return value, true
		}
	}

	return "", false
}

// ReleaseLock 释放分布式锁（使用Lua脚本确保原子性）
func (r *RedisLockService) ReleaseLock(key, value string) bool {
	// Lua脚本原子性地检查并删除锁
	luaScript := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `

	result, err := r.client.Eval(luaScript, []string{key}, value).Result()
	if err != nil {
		return false
	}

	return result.(int64) == 1
}

// IsHeldByCurrent 检查锁是否由当前实例持有
func (r *RedisLockService) IsHeldByCurrent(key, value string) bool {
	val, err := r.client.Get(key).Result()
	if err != nil {
		return false
	}
	return val == value
}

// PaymentRepository 支付记录存储接口
type PaymentRepository interface {
	GetPaymentRecord(idempotentKey string) *PaymentRecord
	CreatePaymentRecord(idempotentKey string, record *PaymentRecord) error
	UpdatePaymentRecord(idempotentKey string, status PaymentStatus, result *PaymentResult, errorMsg string) error
	DeletePaymentRecord(idempotentKey string) error
}

// MemoryPaymentRepository 内存存储实现（模拟数据库）
type MemoryPaymentRepository struct {
	records map[string]*PaymentRecord
	mutex   sync.RWMutex
}

func NewMemoryPaymentRepository() *MemoryPaymentRepository {
	return &MemoryPaymentRepository{
		records: make(map[string]*PaymentRecord),
	}
}

/*
*

		实际根据业务场景来：
		1. 如果只需要短期幂等，那么可以把数据存在redis中，并设置TTL过期时间
		2. 如果需要长期幂等，那么可以把数据存在数据库中
	 *
*/
func (m *MemoryPaymentRepository) GetPaymentRecord(idempotentKey string) *PaymentRecord {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if record, exists := m.records[idempotentKey]; exists {
		// 返回副本避免并发修改
		copy := *record
		if record.Result != nil {
			resultCopy := *record.Result
			copy.Result = &resultCopy
		}
		return &copy
	}
	return nil
}

func (m *MemoryPaymentRepository) CreatePaymentRecord(idempotentKey string, record *PaymentRecord) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.records[idempotentKey]; exists {
		return fmt.Errorf("记录已存在")
	}

	m.records[idempotentKey] = record
	return nil
}

func (m *MemoryPaymentRepository) UpdatePaymentRecord(idempotentKey string, status PaymentStatus, result *PaymentResult, errorMsg string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if record, exists := m.records[idempotentKey]; exists {
		record.Status = status
		record.Result = result
		record.ErrorMessage = errorMsg
		record.UpdatedAt = time.Now()
		return nil
	}
	return fmt.Errorf("记录不存在")
}

func (m *MemoryPaymentRepository) DeletePaymentRecord(idempotentKey string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.records, idempotentKey)
	return nil
}

// PaymentService 支付服务
type PaymentService struct {
	lockService LockService
	repository  PaymentRepository
	requestID   string
}

// LockService 锁服务接口
type LockService interface {
	AcquireLock(key string, expiration time.Duration) (string, bool)
	TryAcquireLock(key string, waitTime, expiration time.Duration) (string, bool)
	ReleaseLock(key, value string) bool
	IsHeldByCurrent(key, value string) bool
}

func NewPaymentService(lockService LockService, repository PaymentRepository) *PaymentService {
	return &PaymentService{
		lockService: lockService,
		repository:  repository,
	}
}

// ProcessPaymentWithIdempotent 处理幂等性支付请求
func (ps *PaymentService) ProcessPaymentWithIdempotent(request PaymentRequest) PaymentResult {
	idempotentKey := request.IdempotentToken

	// 参数校验。通常除了验空之外，还需要验证幂等key的有效性，防篡改等。
	// 这里为了代码简单就跳过对key的有效性校验
	if idempotentKey == "" {
		return PaymentResult{Status: "ERROR", Message: "缺少幂等token"}
	}

	// 1. 先查询是否已有处理结果。
	// 这里可以根据业务场景使用多级缓存，比如 先查本地缓存(L1缓存) -> 查询Redis(L2缓存) -> 再查db等
	if existingRecord := ps.getPaymentRecord(idempotentKey); existingRecord != nil {
		switch existingRecord.Status {
		case StatusSuccess:
			result := *existingRecord.Result
			result.RequestID = ps.requestID
			return result
		case StatusProcessing:
			// 检查是否超时 (通常设置为业务处理超时时间)
			if time.Since(existingRecord.UpdatedAt) > 60*time.Second {
				// 超时，允许重试，更新状态为失败
				ps.updateRecordStatus(idempotentKey, StatusFailed, nil, "处理超时")
			} else {
				return PaymentResult{
					Status:    "PROCESSING",
					Message:   "请求正在处理中",
					RequestID: ps.requestID,
				}
			}
		case StatusFailed:
			// 根据业务决定是否允许重试
			if ps.allowRetry(request) {
				// 清除失败记录，允许重试
				ps.repository.DeletePaymentRecord(idempotentKey)
			} else {
				return PaymentResult{
					Status:    "FAILED",
					Message:   "请求已失败且不可重试",
					RequestID: ps.requestID,
				}
			}
		}
	}

	// 2. 获取分布式锁
	lockKey := "payment_lock:" + idempotentKey
	lockValue, acquired := ps.lockService.TryAcquireLock(lockKey, 3*time.Second, 10*time.Second)
	if !acquired {
		return PaymentResult{
			Status:    "RETRY",
			Message:   "系统繁忙，请稍后重试", // 获取锁失败，可提示用户稍后再试 或 订单正在处理中
			RequestID: ps.requestID,
		}
	}

	// 确保锁被释放
	defer func() {
		if ps.lockService.IsHeldByCurrent(lockKey, lockValue) {
			ps.lockService.ReleaseLock(lockKey, lockValue)
		}
	}()

	// 3. 双重检查：获取锁后再次检查
	if existingRecord := ps.getPaymentRecord(idempotentKey); existingRecord != nil {
		if existingRecord.Status == StatusSuccess {
			result := *existingRecord.Result
			result.RequestID = ps.requestID
			return result
		}
	}

	// 4. 创建处理中记录
	if err := ps.createProcessingRecord(idempotentKey); err != nil {
		return PaymentResult{
			Status:    "ERROR",
			Message:   "创建处理记录失败: " + err.Error(),
			RequestID: ps.requestID,
		}
	}

	// 5. 执行业务逻辑
	var result PaymentResult
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = PaymentResult{
					Status:    "FAILED",
					Message:   "处理过程发生异常",
					RequestID: ps.requestID,
				}
				ps.repository.UpdatePaymentRecord(idempotentKey, StatusFailed, &result, "处理过程发生异常")
			}
		}()

		result = ps.executePayment(request)
		result.RequestID = ps.requestID
		ps.repository.UpdatePaymentRecord(idempotentKey, StatusSuccess, &result, "")
	}()

	return result
}

// getPaymentRecord 获取支付记录
func (ps *PaymentService) getPaymentRecord(idempotentKey string) *PaymentRecord {
	return ps.repository.GetPaymentRecord(idempotentKey)
}

// updateRecordStatus 更新记录状态
func (ps *PaymentService) updateRecordStatus(idempotentKey string, status PaymentStatus, result *PaymentResult, errorMsg string) {
	ps.repository.UpdatePaymentRecord(idempotentKey, status, result, errorMsg)
}

// createProcessingRecord 创建处理中记录
func (ps *PaymentService) createProcessingRecord(idempotentKey string) error {
	record := &PaymentRecord{
		IdempotentKey: idempotentKey,
		Status:        StatusProcessing,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return ps.repository.CreatePaymentRecord(idempotentKey, record)
}

// allowRetry 是否允许重试（根据业务需求实现）
func (ps *PaymentService) allowRetry(request PaymentRequest) bool {
	// 根据业务规则判断是否允许重试，例如检查失败次数等
	return true
}

// executePayment 执行支付逻辑（模拟）
func (ps *PaymentService) executePayment(request PaymentRequest) PaymentResult {
	log.Printf("开始处理支付请求: 用户=%s, 金额=%.2f, 描述=%s",
		request.UserID, request.Amount, request.Description)

	// 模拟支付处理时间（100-500ms）
	n, _ := rand.Int(rand.Reader, big.NewInt(400))
	time.Sleep(time.Duration(100+n.Int64()) * time.Millisecond)

	// 模拟随机失败(5%概率失败)
	failureRate, _ := rand.Int(rand.Reader, big.NewInt(100))
	if failureRate.Int64() < 5 {
		log.Printf("支付处理失败: 用户=%s, 金额=%.2f", request.UserID, request.Amount)
		return PaymentResult{
			Status:  "FAILED",
			Message: "支付网关暂时不可用，请稍后重试",
		}
	}

	orderID := "PAY_" + strings.ToUpper(uuid.New().String()[:8])
	log.Printf("支付处理成功: 用户=%s, 金额=%.2f, 订单=%s",
		request.UserID, request.Amount, orderID)

	return PaymentResult{
		Status:  "SUCCESS",
		Message: fmt.Sprintf("支付成功，金额: %.2f元", request.Amount),
		OrderID: orderID,
	}
}

// HTTP handlers
type PaymentHandler struct {
	paymentService *PaymentService
	repository     PaymentRepository
}

func NewPaymentHandler(paymentService *PaymentService, repository PaymentRepository) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		repository:     repository,
	}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var request PaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 为每次请求生成唯一ID用于追踪
	h.paymentService.requestID = "REQ_" + strings.ToUpper(uuid.New().String()[:8])

	result := h.paymentService.ProcessPaymentWithIdempotent(request)

	// 根据状态码返回不同的HTTP状态
	var statusCode int
	switch result.Status {
	case "SUCCESS":
		statusCode = http.StatusOK
	case "PROCESSING":
		statusCode = http.StatusAccepted
	case "FAILED", "ERROR":
		statusCode = http.StatusInternalServerError
	case "RETRY":
		statusCode = http.StatusTooManyRequests
	default:
		statusCode = http.StatusOK
	}

	c.JSON(statusCode, result)
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	idempotentToken := c.Query("token")
	if idempotentToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少token参数"})
		return
	}

	record := h.repository.GetPaymentRecord(idempotentToken)
	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到对应的支付记录"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *PaymentHandler) ListAllRecords(c *gin.Context) {
	// 这里简单返回所有记录（实际生产中应该分页）
	c.JSON(http.StatusOK, h.repository.(*MemoryPaymentRepository).records)
}

func main() {
	// 初始化Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// 测试Redis连接
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatal("无法连接到Redis: ", err)
	}
	log.Println("成功连接到Redis")

	// 初始化服务
	lockService := NewRedisLockService(redisClient)
	repository := NewMemoryPaymentRepository()
	paymentService := NewPaymentService(lockService, repository)
	paymentHandler := NewPaymentHandler(paymentService, repository)

	// 初始化Gin路由器
	r := gin.Default()

	// API路由
	r.POST("/payment", paymentHandler.ProcessPayment)
	r.GET("/payment/status", paymentHandler.GetPaymentStatus)
	r.GET("/payment/records", paymentHandler.ListAllRecords)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	log.Println("服务器启动在端口 8080")
	log.Println("API endpoints:")
	log.Println("  POST /payment - 处理支付请求")
	log.Println("  GET  /payment/status?token={token} - 查询支付状态")
	log.Println("  GET  /payment/records - 查看所有支付记录")
	log.Println("  GET  /health - 健康检查")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
