package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaorui/reddit-async/reddit-backend/dao/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/logic"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
)

// CreateComment： 创建一个评论
//	@Summary		创建一个评论
//	@Description	创建一个评论
//	@Tags			Comment
//	@Accept			application/json
//	@Produce		application/json
//	@Param			post_id			path	int		true	"Post ID"
//	@Param			Authorization	header	string	false	"Bearer 用户令牌"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/post/comment/{post_id} [post]
func CreateComment(ctx *gin.Context) {
	// 获取要给那个post发送comment
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//1. 进行参数验证
	p := new(models.ParamCreateNewComment)
	if ok := Validate(ctx, p, ValidateCreateComment); !ok {
		return
	}

	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//2. 处理业务逻辑
	if err := logic.CreateComment(postID, userID, p); err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//3. 返回响应
	ResponseSuccess(ctx, nil)
}

// GetAllComment： 获得某个post下的所有评论
//	@Summary		获得某个post下的所有评论
//	@Description	获得某个post下的所有评论
//	@Tags			Comment
//	@Accept			application/json
//	@Produce		application/json
//	@Param			post_id			path	int		true	"Post ID"
//	@Param			Authorization	header	string	false	"Bearer 用户令牌"
//	@Param			page_num		query	int		false	"Page number"
//	@Param			page_size		query	int		false	"Page size"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/post/comment/{post_id} [get]
func GetComment(ctx *gin.Context) {
	//1. 进行参数验证
	postIDStr := ctx.Param("post_id")
	pageNumStr, pageSizeStr := ctx.Query("page_num"), ctx.Query("page_size")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	pageNum, err := strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = 10
	}

	//2. 处理业务逻辑
	commentList, err := logic.GetComment(postID, pageNum, pageSize)
	if err != nil {
		if err == mysql.ErrorCommentNotFound {
			ResponseError(ctx, CodeCommentNotFound)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//3. 返回响应
	ResponseSuccess(ctx, commentList)

}

// DeleteComment： 删除某条评论
//	@Summary		删除某条评论
//	@Description	删除某条评论
//	@Tags			Comment
//	@Accept			application/json
//	@Produce		application/json
//	@Param			comment_id		path	int		true	"comment ID"
//	@Param			Authorization	header	string	false	"Bearer 用户令牌"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/post/comment/{comment_id} [delete]
func DeleteComment(ctx *gin.Context) {
	commentIDStr := ctx.Param("comment_id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	// 获得当下用户的id
	//必须是自己发起的才能够删除
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//1. 进行参数验证
	//2. 处理业务逻辑
	if err := logic.DeleteComment(userID, commentID); err != nil {
		if err == mysql.ErrorCommentNotFound {
			ResponseError(ctx, CodeCommentNotFound)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//3. 返回响应
	ResponseSuccess(ctx, nil)
}
