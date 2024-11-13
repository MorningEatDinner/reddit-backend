package models

import "time"

type Comment struct {
	ID          int64     `json:"comment_id" gorm:"column:comment_id"`
	PostID      int64     `json:"post_id" gorm:"column:post_id"`
	AuthorID    int64     `json:"author_id" gorm:"column:author_id"`
	Content     string    `json:"content" gorm:"column:content"`
	CreateTime  time.Time `json:"-" gorm:"column:create_time;autoCreateTime"`
	UpdatedTime time.Time `json:"-" gorm:"column:updated_time;autoUpdateTime"`
	User        User      `json:"-" gorm:"foreignKey:AuthorID"`
	Post        Post      `json:"-" gorm:"foreignKey:PostID"`
}
