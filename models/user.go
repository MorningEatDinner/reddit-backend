package models

import "time"

type User struct {
	ID           int64     `json:"user_id,string" gorm:"primaryKey;column:user_id"`
	Username     string    `json:"username" gorm:"column:username"`
	Password     string    `json:"-" gorm:"column:password"` //  使用 “- ” 使得在序列化的结果中不会出现当前字段
	Email        string    `json:"email" gorm:"column:email"`
	Phone        string    `json:"phone" gorm:"column:phone"`
	City         string    `json:"city" gorm:"column:city"`
	Introduction string    `json:"introduction" gorm:"column:introduction"`
	Avatar       string    `json:"avatar" gorm:"column:avatar"`
	Token        string    // 注意：此字段没有json或gorm标签
	CreateTime   time.Time `json:"-" gorm:"column:create_time;autoCreateTime"`
	UpdatedTime  time.Time `json:"-" gorm:"column:updated_time;autoUpdateTime"`
}
