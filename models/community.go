package models

import "time"

// type Community struct {
// 	ID   int64  `json:"community_id" db:"community_id"`
// 	Name string `json:"community_name" db:"community_name"`
// }

//	type CommunityDetail struct {
//		ID           int64     `json:"community_id" db:"community_id"`
//		Name         string    `json:"community_name" db:"community_name"`
//		Introduction string    `json:"introduction" db:"introduction"`
//		CreateTime   time.Time `json:"create_time" db:"create_time"`
//	}
type Community struct {
	ID           int64     `json:"community_id" gorm:"column:community_id"`
	Name         string    `json:"community_name" gorm:"column:community_name"`
	Introduction string    `json:"introduction,omitempty" gorm:"column:introduction"`
	CreateTime   time.Time `json:"-" gorm:"column:create_time;autoCreateTime"`
	UpdatedTime  time.Time `json:"-" gorm:"column:updated_time;autoUpdateTime"`
}
