package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/ziyifast/log"
)

func AuditTrailMiddleware(ctx iris.Context) {
	log.Infof("audit trail ....")
	ctx.Next()
}
