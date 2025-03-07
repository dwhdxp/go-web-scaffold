package controllers

import (
	"friend-ranking/cache"
	"friend-ranking/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type PlayerController struct{}

// 获取活动参赛选手详情
func (p PlayerController) GetPlayers(c *gin.Context) {
	aidStr := c.DefaultPostForm("aid", "0")
	aid, _ := strconv.Atoi(aidStr)
	res, err := models.GetPlayers(aid, "id Asc")
	if err != nil {
		ReturnError(c, 4004, "没有相关选手信息")
		return
	}
	ReturnSuccess(c, 0, "success", res, 1)
}

// 排行榜功能
func (p PlayerController) GetRanking(c *gin.Context) {
	aidStr := c.DefaultPostForm("aid", "0")
	aid, _ := strconv.Atoi(aidStr)

	// 从redis中获取排行榜信息
	var redisKey string
	redisKey = "ranking:" + aidStr
	// 反向（降序）获取有序集合Key中排名在[start, stop]范围内的所有数据
	res, err := cache.Rdb.ZRevRange(redisKey, 0, -1).Result()
	if err == nil && len(res) > 0 {
		// 根据player_id查找对应info
		var players []models.Player
		for _, value := range res {
			id, _ := strconv.Atoi(value)
			pInfo, _ := models.GetPlayerInfoById(id)
			players = append(players, pInfo)
		}
		ReturnSuccess(c, 0, "success form redis", players, 1)
		return
	}

	// 从数据库中获取排行榜信息
	resDb, errDb := models.GetPlayers(aid, "score Desc")
	if errDb == nil {
		// 将mysql中获取到的数据缓存到redis
		for _, value := range resDb {
			// 向有序集合中添加一个或多个成员，并为其指定分数
			cache.Rdb.ZAdd(redisKey, cache.Zscore(value.Id, value.Score)).Err()
		}
		// 设置redisKey过期时间
		cache.Rdb.Expire(redisKey, time.Hour)
		ReturnSuccess(c, 0, "success form mysql", resDb, 1)
		return
	}
	ReturnError(c, 4004, "没有相关选手信息")
}
