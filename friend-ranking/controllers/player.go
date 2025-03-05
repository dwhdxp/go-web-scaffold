package controllers

import (
	"friend-ranking/models"
	"github.com/gin-gonic/gin"

	"strconv"
)

type PlayerController struct{}

// 获取活动参赛选手详情
func (p PlayerController) GetPlayers(c *gin.Context) {
	aidStr := c.DefaultPostForm("aid", "0")
	aid, _ := strconv.Atoi(aidStr)
	res, err := models.GetPlayers(aid)
	if err != nil {
		ReturnError(c, 4004, "没有相关选手信息")
		return
	}
	ReturnSuccess(c, 0, "success", res, 1)
}
