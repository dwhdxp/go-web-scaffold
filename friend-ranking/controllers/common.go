package controllers

import (
	"crypto/md5" // 计算Md5哈希
	"encoding/hex"
	"github.com/gin-gonic/gin"
)

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}

type JsonErrStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, data interface{}, count int64) {
	json := &JsonStruct{Code: code, Msg: msg, Data: data, Count: count}
	c.JSON(200, json)
}

func ReturnError(c *gin.Context, code int, msg interface{}) {
	json := &JsonErrStruct{Code: code, Msg: msg}
	c.JSON(200, json)
}

// Md5加密
func EncryMd5(p string) string {
	h := md5.New()
	h.Write([]byte(p))
	hBytes := h.Sum(nil)
	// 将字节数组转换为十六进制字符串
	return hex.EncodeToString(hBytes)
}
