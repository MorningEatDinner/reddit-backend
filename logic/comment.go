package logic

import (
	"github.com/xiaorui/reddit-async/reddit-backend/dao/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/snowflake"
	"go.uber.org/zap"
)

// CreateComment: 给定postid创建一个新的评论
func CreateComment(postID, userID int64, p *models.ParamCreateNewComment) error {
	// 1. 先看post是否存在
	_, err := mysql.GetPostByID(postID) // 如果post存在则不会返回错误
	if err != nil {
		zap.L().Error("mysql.GetPostByID failed...", zap.Error(err))
		return err
	}

	// 生成commentid
	commentID := snowflake.GenID()

	// 创建新的实例
	comment := &models.Comment{
		ID:       commentID,
		AuthorID: userID,
		PostID:   postID,
		Content:  p.Content,
	}

	// 2. 保存内容
	return mysql.CreateComment(comment)
}

// DeleteComment: 删除Comment
func DeleteComment(userid, commentID int64) error {
	//1. 查询Comment是否存在
	comment, err := mysql.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	// 2. 如果Comment存在则删除
	return mysql.DeleteComment(userid, comment)
}

// GetComment: 返回给定post的评论
func GetComment(postID, pageNum, pageSize int64) (comms []*models.Comment, err error) {
	//1. 验证post是否存在
	if _, err = mysql.GetPostByID(postID); err != nil {
		return
	}
	// 2. 查询返回结果
	return mysql.GetComments(postID, pageNum, pageSize)
}
