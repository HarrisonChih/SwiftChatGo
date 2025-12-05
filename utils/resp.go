package utils

import (
	"github.com/goccy/go-json"
	"net/http"
)

// HTTP 接口统一响应工具
// 核心作用是为项目所有 API 接口提供标准化的 JSON 响应格式
type H struct {
	Code  int         // 响应状态码（0=成功，-1=失败，可扩展其他状态）
	Msg   string      // 响应提示信息（如“登录成功”“参数错误”）
	Data  interface{} // 单条数据（如用户信息、详情数据）
	Rows  interface{} // 列表数据（如用户列表、消息列表）
	Total interface{} // 列表总数（用于分页场景）
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h := H{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	w.Write(ret)
}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

func RespOK(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	w.Write(ret)
}

func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}

func RespFailList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, -1, nil, nil)
}
