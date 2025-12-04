package tcc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// TCCTransactionManager TCC事务管理器
type TCCTransactionManager struct {
	participants []TCCParticipant
	mutex        sync.RWMutex
	txID         string
}

// TCCParticipant TCC事务参与者接口
type TCCParticipant interface {
	Try(ctx context.Context) error     // 尝试执行业务
	Confirm(ctx context.Context) error // 确认提交
	Cancel(ctx context.Context) error  // 取消回滚
	GetID() string                     // 获取参与者ID
}

// NewTCCTransactionManager 创建新的TCC事务管理器
func NewTCCTransactionManager(txID string) *TCCTransactionManager {
	return &TCCTransactionManager{
		participants: make([]TCCParticipant, 0),
		txID:         txID,
	}
}

// AddParticipant 添加TCC事务参与者
func (tm *TCCTransactionManager) AddParticipant(participant TCCParticipant) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.participants = append(tm.participants, participant)
}

// ExecuteTCC 执行TCC事务
func (tm *TCCTransactionManager) ExecuteTCC(ctx context.Context) error {
	// 设置默认超时时间
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Try阶段
	if err := tm.tryPhase(ctxWithTimeout); err != nil {
		// Try阶段失败，执行Cancel阶段
		tm.cancelAll(ctxWithTimeout)
		return fmt.Errorf("try phase failed: %w", err)
	}

	// Confirm阶段
	if err := tm.confirmPhase(ctxWithTimeout); err != nil {
		log.Printf("CRITICAL: confirm phase failed for transaction %s: %v", tm.txID, err)
		return fmt.Errorf("confirm phase failed: %w", err)
	}

	return nil
}

// tryPhase Try阶段 - 尝试执行业务
func (tm *TCCTransactionManager) tryPhase(ctx context.Context) error {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	log.Printf("Starting Try phase for transaction %s with %d participants", tm.txID, len(tm.participants))

	for i, participant := range tm.participants {
		if err := participant.Try(ctx); err != nil {
			return fmt.Errorf("participant %s Try failed at index %d: %w", participant.GetID(), i, err)
		}
		log.Printf("Participant %s Try succeeded", participant.GetID())
	}

	log.Printf("Try phase completed successfully for transaction %s", tm.txID)
	return nil
}

// confirmPhase Confirm阶段 - 确认提交
func (tm *TCCTransactionManager) confirmPhase(ctx context.Context) error {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	log.Printf("Starting Confirm phase for transaction %s", tm.txID)

	for i, participant := range tm.participants {
		if err := participant.Confirm(ctx); err != nil {
			log.Printf("CRITICAL: Participant %s Confirm failed at index %d. Manual intervention may be required.", participant.GetID(), i)
			return fmt.Errorf("participant %s Confirm failed at index %d: %w", participant.GetID(), i, err)
		}
		log.Printf("Participant %s Confirm succeeded", participant.GetID())
	}

	log.Printf("Confirm phase completed successfully for transaction %s", tm.txID)
	return nil
}

// cancelAll Cancel阶段 - 全部回滚
func (tm *TCCTransactionManager) cancelAll(ctx context.Context) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	log.Printf("Canceling all participants for transaction %s", tm.txID)

	for i, participant := range tm.participants {
		if err := participant.Cancel(ctx); err != nil {
			log.Printf("Participant %s Cancel failed at index %d: %v", participant.GetID(), i, err)
		} else {
			log.Printf("Participant %s Cancel succeeded", participant.GetID())
		}
	}
}

// AccountServiceParticipant 账户服务参与者示例
type AccountServiceParticipant struct {
	id       string
	userID   int
	amount   float64
	db       interface{} // 实际项目中应替换为具体数据库连接
	reserved bool        // 标记是否已预留资源
}

// NewAccountServiceParticipant 创建账户服务参与者
func NewAccountServiceParticipant(id string, userID int, amount float64, db interface{}) *AccountServiceParticipant {
	return &AccountServiceParticipant{
		id:     id,
		userID: userID,
		amount: amount,
		db:     db,
	}
}

// Try 尝试阶段 - 冻结账户资金
func (asp *AccountServiceParticipant) Try(ctx context.Context) error {
	// 检查账户余额是否足够
	// balance := asp.getAccountBalance(asp.userID)
	// if balance < asp.amount {
	//     return fmt.Errorf("insufficient balance for user %d", asp.userID)
	// }

	// 冻结资金（预留资源）
	// asp.freezeAccountFund(asp.userID, asp.amount)
	asp.reserved = true

	log.Printf("Account service participant %s: Try phase succeeded - frozen amount %.2f for user %d",
		asp.id, asp.amount, asp.userID)
	return nil
}

// Confirm 确认阶段 - 扣除账户资金
func (asp *AccountServiceParticipant) Confirm(ctx context.Context) error {
	if !asp.reserved {
		return fmt.Errorf("resources not reserved for participant %s", asp.id)
	}

	// 实际扣除账户资金
	// asp.deductAccountFund(asp.userID, asp.amount)

	log.Printf("Account service participant %s: Confirm phase succeeded - deducted amount %.2f for user %d",
		asp.id, asp.amount, asp.userID)
	return nil
}

// Cancel 取消阶段 - 解冻账户资金
func (asp *AccountServiceParticipant) Cancel(ctx context.Context) error {
	if !asp.reserved {
		log.Printf("Account service participant %s: No resources to cancel", asp.id)
		return nil
	}

	// 解冻账户资金
	// asp.unfreezeAccountFund(asp.userID, asp.amount)
	asp.reserved = false

	log.Printf("Account service participant %s: Cancel phase succeeded - unfrozen amount %.2f for user %d",
		asp.id, asp.amount, asp.userID)
	return nil
}

// GetID 获取参与者ID
func (asp *AccountServiceParticipant) GetID() string {
	return asp.id
}

// InventoryServiceParticipant 库存服务参与者示例
type InventoryServiceParticipant struct {
	id        string
	productID string
	quantity  int
	db        interface{} // 实际项目中应替换为具体数据库连接
	reserved  bool        // 标记是否已预留资源
}

// NewInventoryServiceParticipant 创建库存服务参与者
func NewInventoryServiceParticipant(id string, productID string, quantity int, db interface{}) *InventoryServiceParticipant {
	return &InventoryServiceParticipant{
		id:        id,
		productID: productID,
		quantity:  quantity,
		db:        db,
	}
}

// Try 尝试阶段 - 预留库存
func (isp *InventoryServiceParticipant) Try(ctx context.Context) error {
	// 检查库存是否充足
	// stock := isp.getProductStock(isp.productID)
	// if stock < isp.quantity {
	//     return fmt.Errorf("insufficient stock for product %s", isp.productID)
	// }

	// 预留库存（冻结库存）
	// isp.reserveProductStock(isp.productID, isp.quantity)
	isp.reserved = true

	log.Printf("Inventory service participant %s: Try phase succeeded - reserved quantity %d for product %s",
		isp.id, isp.quantity, isp.productID)
	return nil
}

// Confirm 确认阶段 - 扣减库存
func (isp *InventoryServiceParticipant) Confirm(ctx context.Context) error {
	if !isp.reserved {
		return fmt.Errorf("resources not reserved for participant %s", isp.id)
	}

	// 实际扣减库存
	// isp.deductProductStock(isp.productID, isp.quantity)

	log.Printf("Inventory service participant %s: Confirm phase succeeded - deducted quantity %d for product %s",
		isp.id, isp.quantity, isp.productID)
	return nil
}

// Cancel 取消阶段 - 释放库存
func (isp *InventoryServiceParticipant) Cancel(ctx context.Context) error {
	if !isp.reserved {
		log.Printf("Inventory service participant %s: No resources to cancel", isp.id)
		return nil
	}

	// 释放库存（解冻库存）
	// isp.releaseProductStock(isp.productID, isp.quantity)
	isp.reserved = false

	log.Printf("Inventory service participant %s: Cancel phase succeeded - released quantity %d for product %s",
		isp.id, isp.quantity, isp.productID)
	return nil
}

// GetID 获取参与者ID
func (isp *InventoryServiceParticipant) GetID() string {
	return isp.id
}

// UsageExample 使用示例
func UsageExample() {
	// 创建TCC事务管理器
	manager := NewTCCTransactionManager("tcc_tx_12345")

	// 创建参与者（示例中使用mock数据库）
	// db1 := getAccountDBConnection()
	// db2 := getInventoryDBConnection()

	// 创建账户服务参与者
	// accountParticipant := NewAccountServiceParticipant("account_svc", 1, 100.0, db1)

	// 创建库存服务参与者
	// inventoryParticipant := NewInventoryServiceParticipant("inventory_svc", "product_001", 2, db2)

	// 添加参与者到事务管理器
	// manager.AddParticipant(accountParticipant)
	// manager.AddParticipant(inventoryParticipant)

	// 执行TCC事务
	// ctx := context.Background()
	// if err := manager.ExecuteTCC(ctx); err != nil {
	//     log.Fatalf("TCC transaction failed: %v", err)
	// }

	// fmt.Println("TCC transaction completed successfully")
	fmt.Println("TCC usage example created:", manager.txID)
}
