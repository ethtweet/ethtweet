package appWeb

const (
	ResponseSuccessCode  = 0
	ResponseFailCode     = 1
	ResponseNotLoginCode = -1
)

type ResponseFormat struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse(code int, msg string, data interface{}) *ResponseFormat {
	if msg == "" {
		if code == ResponseSuccessCode {
			msg = "success"
		} else if code == ResponseFailCode {
			msg = "fail"
		} else if code == ResponseNotLoginCode {
			msg = "user not login"
		}
	}
	return &ResponseFormat{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
