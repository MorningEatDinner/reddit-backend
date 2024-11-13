package mysql

import (
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"gorm.io/gorm"
)

func CreateComment(comment *models.Comment) error {
	return DB.Model(&models.Comment{}).Create(comment).Error
}

// GetCommentByID: 返回Comment
func GetCommentByID(commentID int64) (*models.Comment, error) {
	comment := &models.Comment{}
	err := DB.Model(&models.Comment{}).Where("comment_id = ?", commentID).First(comment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrorCommentNotFound
	}
	return comment, err
}

// DeleteComment: 删除制定评论
func DeleteComment(userID int64, comment *models.Comment) error {
	if userID != comment.AuthorID {
		return ErrorNotPermission
	}

	return DB.Delete(comment).Error
}

// GetComments: 返回评论
func GetComments(postID, pageNum, pageSize int64) ([]*models.Comment, error) {
	commentList := []*models.Comment{}
	err := DB.Model(&models.Comment{}).Where("post_id = ?", postID).
		Limit(int(pageSize)).Offset(int(pageNum - 1)).
		Find(&commentList).Error
	return commentList, err
}
