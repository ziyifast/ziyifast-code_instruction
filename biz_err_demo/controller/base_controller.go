package controller

import (
	"encoding/json"
	"encoding/xml"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/sirupsen/logrus"
	"myTest/demo_home/biz_err_demo/constant"
	"myTest/demo_home/biz_err_demo/error/biz_err"
	"myTest/demo_home/biz_err_demo/response"
	"net/http"
)

type BaseController struct {
	Ctx iris.Context
}

func commonResp(errMsg string, httpCode int, returnCode response.Code, content interface{}) mvc.Response {
	payload := &response.JsonResponse{
		Code:    returnCode,
		Msg:     errMsg,
		Content: content,
	}

	contentDetail, err := json.Marshal(payload)
	if err != nil {
		logrus.Infof("marshal json response error %v", err)
	}

	return mvc.Response{
		Code:        httpCode,
		Content:     contentDetail,
		ContentType: constant.ContentTypeJson,
	}
}

func (c *BaseController) Xml(httpCode int, content interface{}) mvc.Response {
	payload, err := xml.Marshal(content)
	if err != nil {
		logrus.Errorf("marshal xml response error %v", err)
	}
	return c.XmlRaw(httpCode, payload)
}

func (c *BaseController) XmlOK(content interface{}) mvc.Response {
	payload, err := xml.Marshal(content)
	if err != nil {
		logrus.Errorf("marshal xml response error %v", err)
	}
	return c.XmlRaw(http.StatusOK, payload)
}

func (c *BaseController) XmlRaw(httpCode int, content []byte) mvc.Response {
	return mvc.Response{
		Code:        httpCode,
		Content:     content,
		ContentType: constant.ContentTypeXml,
	}
}

func (c *BaseController) JsonBizError(err error) mvc.Response {
	httpStatus, code, msg := biz_err.ErrResponse(err)
	return commonResp(msg, httpStatus, response.Code(code), nil)
}
