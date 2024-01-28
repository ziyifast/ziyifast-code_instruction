package biz_err

import (
	"myTest/demo_home/biz_err_demo/error/zerr"
	"net/http"
	"strings"
)

const (
	Undefined                 = "Undefined"
	OsCreateFileError         = "OsCreateFileError"
	ImageNotSupported         = "ImageNotSupported"
	UsernameOrPasswordInValid = "UsernameOrPasswordInValid"
)

var errorResponseMap = map[string]string{
	OsCreateFileError:         "创建文件失败",
	ImageNotSupported:         "图片格式不支持",
	UsernameOrPasswordInValid: "用户名或密码错误",
}

var errorHttpStatusMap = map[string]int{
	OsCreateFileError:         http.StatusInternalServerError,
	ImageNotSupported:         http.StatusInternalServerError,
	UsernameOrPasswordInValid: http.StatusInternalServerError,
}

func ParseBizErr(err error) (httpStatus int, code, msg string) {
	if err == nil {
		code = Undefined
	}
	vars := make([]string, 0)
	errWrap := new(zerr.ErrWrap)
	var cause error
	if as := zerr.As(err, &errWrap); as {
		code = errWrap.Code()
		cause = errWrap.Cause()
		vars = errWrap.Vars()
	} else {
		code = Undefined
	}
	if code == Undefined {
		var undefinedMsg string
		if err != nil {
			undefinedMsg = err.Error()
		}
		if undefinedMsg == "" || undefinedMsg == ": " {
			undefinedMsg = errorResponseMap[code]
		}
		return errorHttpStatusMap[code], code, undefinedMsg
	}
	if status, ok := errorHttpStatusMap[code]; ok {
		httpStatus = status
	} else {
		httpStatus = http.StatusOK
	}
	if bizMsg, ok := errorResponseMap[code]; ok {
		for _, v := range vars {
			bizMsg = strings.Replace(bizMsg, "%s", v, 1)
		}
		msg = bizMsg
		if cause != nil {
			_, _, causeMsg := ParseBizErr(cause)
			if causeMsg != "" {
				msg += ", " + causeMsg
			} else {
				msg += ", " + errWrap.Error()
			}
		}
	} else {
		msg = errWrap.Error()
	}
	return httpStatus, code, msg
}

func ErrResponse(err error) (httpStatus int, code, msg string) {
	if err == nil {
		code = Undefined
	}
	vars := make([]string, 0)
	errWrap := new(zerr.ErrWrap)
	var cause error
	if as := zerr.As(err, &errWrap); as {
		code = errWrap.Code()
		cause = errWrap.Cause()
		vars = errWrap.Vars()
	} else {
		code = Undefined
	}
	if status, ok := errorHttpStatusMap[code]; ok {
		httpStatus = status
	} else {
		httpStatus = http.StatusOK
	}
	if bizMsg, ok := errorResponseMap[code]; ok {
		for _, v := range vars {
			bizMsg = strings.Replace(bizMsg, "%s", v, 1)
		}
		msg = bizMsg
		if cause != nil {
			_, _, causeMsg := ErrResponse(cause)
			if causeMsg != "" {
				msg += causeMsg
			} else {
				msg += errWrap.Error()
			}
		}
	} else {
		msg = errWrap.Error()
	}
	return httpStatus, code, msg
}
