package response

import (
	"encoding/json"
	"github.com/kataras/iris/v12/mvc"
	"github.com/sirupsen/logrus"
	"myTest/demo_home/biz_err_demo/constant"
	"myTest/demo_home/biz_err_demo/error/biz_err"
)

type Code string

type JsonResponse struct {
	Code    Code        `json:"code"`
	Msg     string      `json:"msg"`
	Content interface{} `json:"content,omitempty"`
}

func JsonBizError(err error) mvc.Response {
	httpStatus, code, msg := biz_err.ErrResponse(err)
	return commonResp(msg, httpStatus, Code(code), nil)
}

func commonResp(errMsg string, httpCode int, returnCode Code, content interface{}) mvc.Response {
	payload := &JsonResponse{
		Code:    returnCode,
		Msg:     errMsg,
		Content: content,
	}

	contentDetail, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("%v", err)
	}

	return mvc.Response{
		Code:        httpCode,
		Content:     contentDetail,
		ContentType: constant.ContentTypeJson,
	}
}
