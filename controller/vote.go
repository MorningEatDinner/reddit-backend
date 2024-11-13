package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xiaorui/reddit-async/reddit-backend/logic"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"go.uber.org/zap"
)

// PostVoteHandler： 对于某个帖子进行投票
//	@Summary		对于某个帖子进行投票
//	@Description	对于某个帖子进行投票
//	@Tags			Post
//	@Accept			application/json
//	@Produce		application/json
//	@Param			object			body	models.ParamVoteData	false	"参数"
//	@Param			Authorization	header	string					false	"Bearer 用户令牌"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/post/vote [post]
func PostVoteHandler(ctx *gin.Context) {
	//1. 获取参数进行参数校验
	p := new(models.ParamVoteData)
	if err := ctx.ShouldBindJSON(p); err != nil {
		//这里还希望将错误进行转换为中文
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) // 如果存在错误转为中文进行输出
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errData)
		return
	}
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}

	//2. 处理业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("PostVoteHandler logic.VoteForPost failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//3. 返回响应
	ResponseSuccess(ctx, nil)
}
