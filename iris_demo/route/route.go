package route

import (
	"github.com/kataras/iris/v12/context"
)

var ControllerList = make([]*ControllerRoute, 0)

type ControllerRoute struct {
	RouteName       string
	ControllerObj   interface{}
	ServiceSlice    []interface{}
	MiddlewareSlice []context.Handler
	DoneHandleSlice []context.Handler
}

func NewControllerRoute(routeName string, controller interface{}, middlewares []context.Handler, doneHandle []context.Handler) {
	route := &ControllerRoute{
		RouteName:       routeName,
		ControllerObj:   controller,
		MiddlewareSlice: middlewares,
		DoneHandleSlice: doneHandle,
	}
	ControllerList = append(ControllerList, route)
}
