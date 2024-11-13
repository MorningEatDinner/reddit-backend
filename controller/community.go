package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaorui/reddit-async/reddit-backend/dao/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/logic"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"go.uber.org/zap"
)

// CommunityHandler: 获取所有社区的信息
//	@Summary		获取所有社区的信息
//	@Description	获取所有社区的信息
//	@Tags			Community
//	@Accept			application/json
//	@Produce		application/json
//	@Param			Authorization	header	string	false	"Bearer 用户令牌"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/community [get]
func CommunityHandler(ctx *gin.Context) {
	// 查询到所有的数据(id, community_name)
	data, err := logic.GetCommunityList() // 从表中获取数据
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed.", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, data)
}

// CommunityDetailHandler: 获取单个社区的详细信息
//	@Summary		获取单个社区的详细信息
//	@Description	获取单个社区的详细信息
//	@Tags			Community
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id				path	int		true	"Community ID"
//	@Param			Authorization	header	string	false	"Bearer 用户令牌"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/community/{id} [get]
func CommunityDetailHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	communityID, err := strconv.ParseInt(idStr, 10, 64) // 10进制， int64类型
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetail(communityID)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, data)
}

// CreateNewCommunity: 创建新的社区
//	@Summary		创建新的社区
//	@Description	创建新的社区
//	@Tags			Community
//	@Accept			application/json
//	@Produce		application/json
//	@Param			object			body	models.ParamCommunity	true	"参数"
//	@Param			Authorization	header	string					false	"Bearer 用户令牌"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/community [post]
func CreateNewCommunity(ctx *gin.Context) {
	p := new(models.ParamCommunity)
	if ok := Validate(ctx, p, ValidateCommunity); !ok {
		return
	}

	// 处理业务： 创建新的社区
	if err := logic.CreateNewCommunity(p); err != nil {
		if err == mysql.ErrorCommunityExist {
			ResponseError(ctx, CodeCommunityExist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccess(ctx, nil)
}

// UpdateCommunity: 更新某个社区的信息
//	@Summary		更新某个社区的信息
//	@Description	更新某个社区的信息
//	@Tags			Community
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id				path	int		true	"Community ID"
//	@Param			Authorization	header	string	false	"Bearer 用户令牌"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/community/{id} [put]
func UpdateCommunity(ctx *gin.Context) {
	// 1. 获取id
	communityID := ctx.Param("id")
	if communityID == "" {
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	// 2. 验证参数
	p := new(models.ParamCommunity)
	if ok := Validate(ctx, p, ValidateCommunity); !ok {
		return
	}

	// 3. 处理逻辑；更新社区信息
	if community, err := logic.UpdateCommunity(communityID, p); err != nil {
		if err == mysql.ErrorCommunityNotExist {
			ResponseError(ctx, CodeCommunityNotEXist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	} else {
		ResponseSuccess(ctx, community)
	}
}

// DeleteCommunity： 删除某个社区的信息
//	@Summary		删除某个社区的信息
//	@Description	删除某个社区的信息
//	@Tags			Community
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id				path	int		true	"Community ID"
//	@Param			Authorization	header	string	false	"Bearer 用户令牌"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/community/{id} [delete]
func DeleteCommunity(ctx *gin.Context) {
	// 1. 获取id
	communityID := ctx.Param("id")
	if communityID == "" {
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	// 2. 进行删除业务
	if err := logic.DeleteCommunity(communityID); err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccess(ctx, nil)
}
