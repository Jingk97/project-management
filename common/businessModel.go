package common

// 定义标准业务数据模型

type BusinessCode int
type Result struct {
	Code BusinessCode
	Msg  string
	Data any
}

func (r *Result) Success(data any) *Result {
	r.Code = 200
	r.Msg = "success"
	r.Data = data
	return r
}

func (r *Result) Fail(code BusinessCode, msg string) *Result {
	r.Code = code
	r.Msg = msg
	return r
}
