package at

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// ATTransactionManager AT事务管理器
type ATTransactionManager struct {
	participants []ATParticipant
	mutex        sync.RWMutex
	txID         string
}

// ATParticipant AT事务参与者接口
type ATParticipant interface {
	ExecuteBusinessLogic(ctx context.Context) error // 执行业务逻辑并记录Undo日志
	Commit(ctx context.Context) error               // 提交事务
	Rollback(ctx context.Context) error             // 回滚事务
	GetID() string                                  // 获取参与者ID
}

// NewATTransactionManager 创建新的AT事务管理器
func NewATTransactionManager(txID string) *ATTransactionManager {
	return &ATTransactionManager{
		participants: make([]ATParticipant, 0),
		txID:         txID,
	}
}

// AddParticipant 添加AT事务参与者
func (tm *ATTransactionManager) AddParticipant(participant ATParticipant) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.participants = append(tm.participants, participant)
}

// ExecuteAT 执行AT事务
func (tm *ATTransactionManager) ExecuteAT(ctx context.Context) error {
	// 设置默认超时时间
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 执行业务逻辑阶段
	if err := tm.businessPhase(ctxWithTimeout); err != nil {
		// 业务执行失败，执行回滚阶段
		tm.rollbackAll(ctxWithTimeout)
		return fmt.Errorf("business phase failed: %w", err)
	}

	// 提交阶段
	if err := tm.commitPhase(ctxWithTimeout); err != nil {
		log.Printf("CRITICAL: commit phase failed for transaction %s: %v", tm.txID, err)
		return fmt.Errorf("commit phase failed: %w", err)
	}

	return nil
}

// businessPhase 业务执行阶段 - 执行业务逻辑并记录Undo日志
func (tm *ATTransactionManager) businessPhase(ctx context.Context) error {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	log.Printf("Starting business phase for transaction %s with %d participants", tm.txID, len(tm.participants))

	for i, participant := range tm.participants {
		if err := participant.ExecuteBusinessLogic(ctx); err != nil {
			return fmt.Errorf("participant %s business logic execution failed at index %d: %w", participant.GetID(), i, err)
		}
		log.Printf("Participant %s business logic executed successfully", participant.GetID())
	}

	log.Printf("Business phase completed successfully for transaction %s", tm.txID)
	return nil
}

// commitPhase 提交阶段 - 提交所有参与者
func (tm *ATTransactionManager) commitPhase(ctx context.Context) error {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	log.Printf("Starting commit phase for transaction %s", tm.txID)

	for i, participant := range tm.participants {
		if err := participant.Commit(ctx); err != nil {
			log.Printf("CRITICAL: Participant %s commit failed at index %d. Manual intervention may be required.", participant.GetID(), i)
			return fmt.Errorf("participant %s commit failed at index %d: %w", participant.GetID(), i, err)
		}
		log.Printf("Participant %s committed successfully", participant.GetID())
	}

	log.Printf("Commit phase completed successfully for transaction %s", tm.txID)
	return nil
}

// rollbackAll 回滚阶段 - 回滚所有参与者
func (tm *ATTransactionManager) rollbackAll(ctx context.Context) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	log.Printf("Rolling back all participants for transaction %s", tm.txID)

	for i, participant := range tm.participants {
		if err := participant.Rollback(ctx); err != nil {
			log.Printf("Participant %s rollback failed at index %d: %v", participant.GetID(), i, err)
		} else {
			log.Printf("Participant %s rolled back successfully", participant.GetID())
		}
	}
}

// UndoLog Undo日志结构
type UndoLog struct {
	ID           int64       `json:"id"`
	BranchID     string      `json:"branch_id"`
	XID          string      `json:"xid"`
	Context      string      `json:"context"`
	RollbackInfo interface{} `json:"rollback_info"`
	LogStatus    int         `json:"log_status"`
	LogCreated   time.Time   `json:"log_created"`
	LogModified  time.Time   `json:"log_modified"`
}

// BusinessOperation 业务操作函数类型
type BusinessOperation func(*sql.Tx) ([]UndoRecord, error)

// UndoRecord 回滚记录
type UndoRecord struct {
	TableName string                 `json:"table_name"`
	SQLType   string                 `json:"sql_type"` // INSERT, UPDATE, DELETE
	Before    map[string]interface{} `json:"before"`   // 操作前数据
	After     map[string]interface{} `json:"after"`    // 操作后数据
}

// DatabaseParticipant 数据库参与者
type DatabaseParticipant struct {
	db         *sql.DB
	txID       string
	branchID   string
	id         string
	operations []BusinessOperation
	undoLogs   []UndoRecord
	currentTx  *sql.Tx
}

// NewDatabaseParticipant 创建新的数据库参与者
func NewDatabaseParticipant(db *sql.DB, txID, branchID, id string) *DatabaseParticipant {
	return &DatabaseParticipant{
		db:         db,
		txID:       txID,
		branchID:   branchID,
		id:         id,
		operations: make([]BusinessOperation, 0),
		undoLogs:   make([]UndoRecord, 0),
	}
}

// AddOperation 添加业务操作
func (dp *DatabaseParticipant) AddOperation(op BusinessOperation) {
	dp.operations = append(dp.operations, op)
}

// ExecuteBusinessLogic 执行业务逻辑并记录Undo日志
func (dp *DatabaseParticipant) ExecuteBusinessLogic(ctx context.Context) error {
	// 创建数据库事务
	tx, err := dp.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	dp.currentTx = tx

	// 执行所有业务操作并收集Undo日志
	for i, op := range dp.operations {
		undoRecords, err := op(tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("business operation %d failed: %w", i, err)
		}

		// 收集Undo日志
		dp.undoLogs = append(dp.undoLogs, undoRecords...)
	}

	log.Printf("Business logic executed successfully for participant %s, collected %d undo records", dp.id, len(dp.undoLogs))
	return nil
}

// Commit 提交事务
func (dp *DatabaseParticipant) Commit(ctx context.Context) error {
	if dp.currentTx != nil {
		if err := dp.currentTx.Commit(); err != nil {
			log.Printf("Transaction commit failed for participant %s: %v", dp.id, err)
			return fmt.Errorf("transaction commit failed: %w", err)
		}

		// 提交成功后可以异步删除Undo日志
		dp.currentTx = nil
		log.Printf("Transaction committed successfully for participant %s", dp.id)
	}

	return nil
}

// Rollback 回滚事务
func (dp *DatabaseParticipant) Rollback(ctx context.Context) error {
	if dp.currentTx != nil {
		// 先回滚数据库事务
		if err := dp.currentTx.Rollback(); err != nil {
			log.Printf("Database transaction rollback failed for participant %s: %v", dp.id, err)
		} else {
			log.Printf("Database transaction rolled back successfully for participant %s", dp.id)
		}
		dp.currentTx = nil
	}

	// 执行Undo操作
	if err := dp.executeUndoLogs(ctx); err != nil {
		log.Printf("Undo logs execution failed for participant %s: %v", dp.id, err)
		return fmt.Errorf("undo logs execution failed: %w", err)
	}

	log.Printf("Participant %s rolled back successfully with %d undo operations", dp.id, len(dp.undoLogs))
	return nil
}

// executeUndoLogs 执行Undo日志回滚
func (dp *DatabaseParticipant) executeUndoLogs(ctx context.Context) error {
	// 反向执行Undo日志
	for i := len(dp.undoLogs) - 1; i >= 0; i-- {
		record := dp.undoLogs[i]

		var err error
		switch record.SQLType {
		case "INSERT":
			// 对于INSERT操作，回滚需要DELETE
			err = dp.executeDeleteUndo(ctx, record)
		case "UPDATE":
			// 对于UPDATE操作，回滚需要恢复到Before状态
			err = dp.executeUpdateUndo(ctx, record)
		case "DELETE":
			// 对于DELETE操作，回滚需要INSERT
			err = dp.executeInsertUndo(ctx, record)
		default:
			err = fmt.Errorf("unsupported SQL type for undo: %s", record.SQLType)
		}

		if err != nil {
			return fmt.Errorf("failed to execute undo log %d: %w", i, err)
		}
	}

	return nil
}

// executeDeleteUndo 执行DELETE类型的Undo操作
func (dp *DatabaseParticipant) executeDeleteUndo(ctx context.Context, record UndoRecord) error {
	// 构造DELETE语句删除新插入的记录
	// 这里简化处理，实际需要根据主键等条件构造WHERE子句
	query := fmt.Sprintf("DELETE FROM %s WHERE /* conditions based on inserted data */", record.TableName)
	log.Printf("Executing undo DELETE: %s", query)

	_, err := dp.db.ExecContext(ctx, query)
	return err
}

// executeUpdateUndo 执行UPDATE类型的Undo操作
func (dp *DatabaseParticipant) executeUpdateUndo(ctx context.Context, record UndoRecord) error {
	// 构造UPDATE语句恢复到Before状态
	// 这里简化处理，实际需要根据主键等条件构造WHERE子句和SET子句
	query := fmt.Sprintf("UPDATE %s SET /* restore before values */ WHERE /* conditions */", record.TableName)
	log.Printf("Executing undo UPDATE: %s", query)

	_, err := dp.db.ExecContext(ctx, query)
	return err
}

// executeInsertUndo 执行INSERT类型的Undo操作
func (dp *DatabaseParticipant) executeInsertUndo(ctx context.Context, record UndoRecord) error {
	// 构造INSERT语句重新插入被删除的记录
	// 这里简化处理，实际需要根据Before数据构造完整的INSERT语句
	query := fmt.Sprintf("INSERT INTO %s /* columns and values from before data */", record.TableName)
	log.Printf("Executing undo INSERT: %s", query)

	_, err := dp.db.ExecContext(ctx, query)
	return err
}

// GetID 获取参与者ID
func (dp *DatabaseParticipant) GetID() string {
	return dp.id
}

// saveUndoLog 保存Undo日志到数据库
func (dp *DatabaseParticipant) saveUndoLog(ctx context.Context, undoLog *UndoLog) error {
	// 实际实现中需要将Undo日志保存到专门的undo_log表中
	query := "INSERT INTO undo_log (branch_id, xid, context, rollback_info, log_status, log_created, log_modified) VALUES (?, ?, ?, ?, ?, ?, ?)"
	log.Printf("Saving undo log: %s", query)

	// _, err := dp.db.ExecContext(ctx, query,
	//     undoLog.BranchID, undoLog.XID, undoLog.Context,
	//     undoLog.RollbackInfo, undoLog.LogStatus,
	//     undoLog.LogCreated, undoLog.LogModified)

	return nil
}

// UsageExample 使用示例
func UsageExample() {
	// 创建AT事务管理器
	manager := NewATTransactionManager("at_tx_12345")

	// 创建数据库连接（示例）
	// db1, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db1")
	// db2, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/db2")

	// 创建参与者
	// participant1 := NewDatabaseParticipant(db1, "at_tx_12345", "branch_001", "account_service")
	// participant2 := NewDatabaseParticipant(db2, "at_tx_12345", "branch_002", "inventory_service")

	// 添加业务操作
	// participant1.AddOperation(func(tx *sql.Tx) ([]UndoRecord, error) {
	//     // 执行业务SQL
	//     _, err := tx.Exec("UPDATE account SET balance = balance - 100 WHERE user_id = 1")
	//     if err != nil {
	//         return nil, err
	//     }
	//
	//     // 记录Undo日志
	//     undoRecords := []UndoRecord{
	//         {
	//             TableName: "account",
	//             SQLType:   "UPDATE",
	//             Before:    map[string]interface{}{"user_id": 1, "balance": 1000},
	//             After:     map[string]interface{}{"user_id": 1, "balance": 900},
	//         },
	//     }
	//     return undoRecords, nil
	// })
	//
	// participant2.AddOperation(func(tx *sql.Tx) ([]UndoRecord, error) {
	//     // 执行业务SQL
	//     _, err := tx.Exec("UPDATE inventory SET stock = stock - 1 WHERE product_id = 1")
	//     if err != nil {
	//         return nil, err
	//     }
	//
	//     // 记录Undo日志
	//     undoRecords := []UndoRecord{
	//         {
	//             TableName: "inventory",
	//             SQLType:   "UPDATE",
	//             Before:    map[string]interface{}{"product_id": 1, "stock": 50},
	//             After:     map[string]interface{}{"product_id": 1, "stock": 49},
	//         },
	//     }
	//     return undoRecords, nil
	// })

	// 添加参与者到事务管理器
	// manager.AddParticipant(participant1)
	// manager.AddParticipant(participant2)

	// 执行AT事务
	// ctx := context.Background()
	// if err := manager.ExecuteAT(ctx); err != nil {
	//     log.Fatalf("AT transaction failed: %v", err)
	// }

	// fmt.Println("AT transaction completed successfully")
	fmt.Println("AT usage example created:", manager.txID)
}
