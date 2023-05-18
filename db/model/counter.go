package model

import "time"

// CounterModel 计数器模型
type CounterModel struct {
	Id        int32     `gorm:"column:id" json:"id"`
	Count     int32     `gorm:"column:count" json:"count"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

type UserMaxScore struct {
	Id          int32     `gorm:"column:id" json:"id"`
	UserId      string    `gorm:"column:user_id" json:"user_id"`
	Score       int32     `gorm:"column:score" json:"score"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
	UpdatedTime time.Time `gorm:"column:update_time" json:"updated_time"`
}
