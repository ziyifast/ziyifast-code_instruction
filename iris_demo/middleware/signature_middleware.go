package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/ziyifast/log"
)

func SignatureMiddleware(ctx iris.Context) {
	log.Infof("SignatureMiddleware...do some signature check")
	ctx.Next()
}
