package models

import "time"

type WarpLog struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UID       string    `gorm:"index" json:"uid"`
	ItemID    string    `json:"item_id"`
	Name      string    `json:"name"`
	ItemType  string    `json:"item_type"`  // Character or Light Cone
	RankType  int       `json:"rank_type"`  // 3, 4, 5
	GachaType string    `json:"gacha_type"` // 11=Starter, 1=LightCone, 2=Standard, 12=CharacterEvent
	Time      time.Time `json:"time"`
}
