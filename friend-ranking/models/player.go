package models

import (
	"friend-ranking/dao"
	"github.com/jinzhu/gorm"
)

type Player struct {
	Id          int    `json:"id"`
	Aid         int    `json:"aid"`
	Ref         string `json:"ref"`
	Nickname    string `json:"nick_name"`
	Declaration string `json:"declaration"`
	Avatar      string `json:"avatar"`
	Score       int    `json:"score"`
	//AddTime     int64  `json:"addTime"`
	//UpdateTime  int64  `json:"updateTime"`
}

func (Player) TableName() string { return "player" }

// 根据活动id获得所有选手信息
func GetPlayers(aid int) ([]Player, error) {
	var players []Player
	err := dao.Db.Where("aid = ?", aid).Find(&players).Error
	return players, err
}

// 根据选手id获得选手信息
func GetPlayerInfoById(id int) (Player, error) {
	var player Player
	err := dao.Db.Where("id = ?", id).First(&player).Error
	return player, err
}

func UpdatePlayerScore(id int) error {
	var player Player
	err := dao.Db.Model(&player).Where("id = ?", id).
		UpdateColumn("score", gorm.Expr("score + ?", 1)).Error
	return err
}
