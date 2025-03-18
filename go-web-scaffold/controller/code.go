package controller

/**
 * 封装业务状态码
 **/

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParams
	CodeServerBusy
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:       "success",
	CodeInvalidParams: "请求参数错误",
	CodeServerBusy:    "服务繁忙",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
