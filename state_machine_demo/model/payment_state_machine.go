package model

type PaymentStatus string

const (
	INIT   PaymentStatus = "INIT"
	PAYING PaymentStatus = "PAYING"
	PAID   PaymentStatus = "PAID"
	FAILED PaymentStatus = "FAILED"
)

type PaymentEvent string

const (
	PAY_CREATE  PaymentEvent = "PAY_CREATE"
	PAY_PROCESS PaymentEvent = "PAY_PROCESS"
	PAY_SUCCESS PaymentEvent = "PAY_SUCCESS"
	PAY_FAIL    PaymentEvent = "PAY_FAIL"
)

var PaymentStateMachine = StateMachine{statusEventMap: map[StatusEventPair]BaseStatus{}}

func init() {
	//支付状态机初始化，包含所有可能的情况
	PaymentStateMachine.accept(nil, PAY_CREATE, INIT)
	PaymentStateMachine.accept(INIT, PAY_PROCESS, PAYING)
	PaymentStateMachine.accept(PAYING, PAY_SUCCESS, PAID)
	PaymentStateMachine.accept(PAYING, PAY_FAIL, FAILED)
}

func GetTargetStatus(sourceStatus PaymentStatus, event PaymentEvent) PaymentStatus {
	status := PaymentStateMachine.getTargetStatus(sourceStatus, event)
	if status != nil {
		return status.(PaymentStatus)
	}
	panic("获取目标状态失败")
}

type PaymentModel struct {
	lastStatus    PaymentStatus
	CurrentStatus PaymentStatus
}

func (pm *PaymentModel) TransferStatusByEvent(event PaymentEvent) {
	targetStatus := GetTargetStatus(pm.CurrentStatus, event)
	if targetStatus != "" {
		pm.lastStatus = pm.CurrentStatus
		pm.CurrentStatus = targetStatus
	} else {
		// 处理异常
		panic("状态转换失败")
	}
}
