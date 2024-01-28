package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"myTest/demo_home/biz_err_demo/error/biz_err"
	"myTest/demo_home/biz_err_demo/error/zerr"
)

func init() {
	logrus.SetReportCaller(true) // 设置日志是否记录被调用的位置，默认值为 false
}

func main() {
	TestWithNoSourceErr()
	TestWithSourceErr()
	TestParseBizErr()
}

func TestWithNoSourceErr() {
	err := zerr.DefaultBizWrap(biz_err.UsernameOrPasswordInValid, "")
	logrus.Errorf("TestWithNoSourceErr %+v", err)
}

func TestWithSourceErr() {
	err := errors.New("invalid image")
	err = zerr.BizWrap(err, biz_err.ImageNotSupported, "")
	logrus.Errorf("TestWithSourceErr %+v", err)
}

func TestParseBizErr() {
	err := errors.New("")
	err = zerr.BizWrap(err, biz_err.ImageNotSupported, "")
	httpStatus, bizCode, msg := biz_err.ParseBizErr(err)
	logrus.Errorf("httpStatus:%d bizCode:%s msg:%s", httpStatus, bizCode, msg)
}
