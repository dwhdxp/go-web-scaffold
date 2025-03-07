package controllers

import (
	"friend-ranking/cache"
	"friend-ranking/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type VoteController struct{}

// 投票功能
func (v VoteController) AddVote(c *gin.Context) {
	// 获取form参数
	userIdStr := c.DefaultPostForm("userId", "0")
	playerIdStr := c.DefaultPostForm("playerId", "0")
	userId, _ := strconv.Atoi(userIdStr)
	playerId, _ := strconv.Atoi(playerIdStr)

	// 参数校验
	// 1.输入项是否完整
	if userId == 0 || playerId == 0 {
		ReturnError(c, 4001, "信息输入错误")
		return
	}

	// 2.用户是否存在
	user, _ := models.GetUserInfoById(userId)
	if user.Id == 0 {
		ReturnError(c, 4001, "用户不存在")
		return
	}

	// 3.选手是否存在
	player, _ := models.GetPlayerInfoById(playerId)
	if player.Id == 0 {
		ReturnError(c, 4001, "参赛选手不存在")
		return
	}

	// 4.用户是否已投票
	vote, _ := models.GetVoteInfo(userId, playerId)
	if vote.Id != 0 {
		ReturnError(c, 4001, "用户已投票，不可重复投票")
		return
	}

	// 投票
	// 投票记录插入数据库
	res, err := models.AddVote(userId, playerId)
	if err == nil {
		// 更新mysql
		if err := models.UpdatePlayerScore(playerId); err != nil {
			ReturnError(c, 4004, "更新mysql中选手分数失败")
			return
		}
		// 更新redis，避免缓存不一致问题
		var redisKey string
		redisKey = "ranking:" + strconv.Itoa(player.Aid)
		cache.Rdb.ZIncrBy(redisKey, 1, playerIdStr)

		ReturnSuccess(c, 0, "投票成功", res, 1)
		return
	}
	ReturnError(c, 4004, "投票失败，请联系工作人员")
}
