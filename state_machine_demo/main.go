package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/ziyifast/log"
	"myTest/demo_home/state_machine_demo/model"
	"time"
)

var (
	testOrder = new(model.PaymentModel)
)

func main() {
	application := iris.New()
	application.Get("/order/create", createOrder)
	application.Get("/order/pay", payOrder)
	application.Get("/order/status", getOrderStatus)
	application.Listen(":8899", nil)
}

func createOrder(context *context.Context) {
	testOrder.CurrentStatus = model.INIT
	context.WriteString("create order succ...")
}

func payOrder(context *context.Context) {
	testOrder.TransferStatusByEvent(model.PAY_PROCESS)
	log.Infof("call third api....")
	//调用第三方支付接口和其他业务处理逻辑
	time.Sleep(time.Second * 15)
	log.Infof("done...")
	testOrder.TransferStatusByEvent(model.PAY_SUCCESS)
}

func getOrderStatus(context *context.Context) {
	context.WriteString(string(testOrder.CurrentStatus))
}
