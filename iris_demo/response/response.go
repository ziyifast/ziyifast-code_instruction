package response

type Code string

var (
	ReturnCodeError   Code = "0"
	ReturnCodeSuccess Code = "1"
	SuccessMsg             = "success"
	SuccessStatus          = "success"
	FailedStatus           = "failed"

	ContentTypeJson = "application/json"
)

type JsonResponse struct {
	Code    Code        `json:"code"`
	Msg     string      `json:"msg"`
	Content interface{} `json:"data"`
}
