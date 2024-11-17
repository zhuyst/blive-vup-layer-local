package main

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	CodeOK            = 0
	CodeBadRequest    = 400
	CodeInternalError = 500
)

func BuildResultOk(data interface{}) *Result {
	return &Result{
		Code: CodeOK,
		Msg:  "ok",
		Data: data,
	}
}

func BuildResultError(code int, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
	}
}
