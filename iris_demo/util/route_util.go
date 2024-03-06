package util

import (
	"github.com/kataras/iris/v12/context"
	"myTest/demo_home/iris_demo/middleware"
	"myTest/demo_home/iris_demo/route"
)

func NewSignRoute(prefix, routeName string, controllerVar interface{}, extraHandler ...context.Handler) {
	middleWares := []context.Handler{middleware.SignatureMiddleware}
	if len(extraHandler) > 0 {
		middleWares = append(middleWares, extraHandler...)
	}
	doneHandler := []context.Handler{middleware.AuditTrailMiddleware}
	route.NewControllerRoute(prefix+routeName, controllerVar, middleWares, doneHandler)
}
