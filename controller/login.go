package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xiaorui/reddit-async/reddit-backend/dao/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/logic"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/jwt"
	"go.uber.org/zap"
)

func LoginHandler(ctx *gin.Context) {
	//1. 进行参数校验
	p := new(models.ParamLogin)
	if err := ctx.ShouldBindJSON(p); err != nil {
		//如果发生获取参数发生了错误
		zap.L().Error("Login with invalid param", zap.Error(err))
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			//如果不是验证器错误
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		//如果是验证器错误, 就是这里捕获错误之后进行翻译返回
		//ctx.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(err.Translate(trans)),
		//})
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(err.Translate(trans)))
	}

	//2. 业务处理
	user, err := logic.Login(p) // 登录之后获取一个token
	if err != nil {
		zap.L().Error("Login with invalid data...", zap.String("username", p.Username), zap.Error(err))
		//ctx.JSON(http.StatusOK, gin.H{
		//	"msg": "登录失败",
		//})
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(ctx, CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorPasswordInvalid) {
			ResponseError(ctx, CodeInvalidPassword)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}

	//3. 返回响应
	ResponseSuccess(ctx, gin.H{
		"user_id":   fmt.Sprintf("%d", user.ID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}

// LoginUsingPhone: 实现使用手机号码进行登陆的功能
//	@Summary		实现使用手机号码进行登陆的功能
//	@Description	实现使用手机号码进行登陆的功能
//	@Tags			Auth
//	@Accept			application/json
//	@Produce		application/json
//	@Param			object	body	models.ParamLoginUsingPhoneWithCode	false	"查询参数"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/auth/login/phone [post]
func LoginUsingPhone(ctx *gin.Context) {
	// 1. 进行参数的验证
	p := new(models.ParamLoginUsingPhoneWithCode)
	if ok := Validate(ctx, p, ValidateLoginUsingPhoneWithCode); !ok {
		return
	}

	// 2. 进行登录操作
	user, err := logic.LoginUsingPhoneWithCode(p)
	if err != nil {
		if errors.Is(err, mysql.ErrorPhoneNotExist) {
			ResponseError(ctx, CodePhoneNotExist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 生成token
	accessToken, refreshToken, err := jwt.GenToken(user.ID, user.Username)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3. 返回执行响应
	ResponseSuccess(ctx, gin.H{
		"user_id":       user.ID,
		"username":      user.Username,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// LoginUsingEmail: 使用邮箱+密码的方式进行登陆
//	@Summary		使用邮箱+密码的方式进行登陆
//	@Description	使用邮箱+密码的方式进行登陆
//	@Tags			Auth
//	@Accept			application/json
//	@Produce		application/json
//	@Param			object	body	models.ParamLoginUsingEmail	false	"查询参数"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/auth/login/email [post]
func LoginUsingEmail(ctx *gin.Context) {
	p := new(models.ParamLoginUsingEmail)
	if ok := Validate(ctx, p, ValidateLoginUsingEmail); !ok {
		return
	}

	// 2. 进行登录操作
	user, err := logic.LoginUsingEmail(p)
	if err != nil {
		if errors.Is(err, mysql.ErrorEmailNotExist) {
			ResponseError(ctx, CodeEmailNotExist)
			return
		} else if errors.Is(err, mysql.ErrorPasswordInvalid) {
			ResponseError(ctx, CodeInvalidPassword)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 生成token
	accessToken, refreshToken, err := jwt.GenToken(user.ID, user.Username)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3. 返回执行响应
	ResponseSuccess(ctx, gin.H{
		"user_id":       user.ID,
		"username":      user.Username,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken: 刷新token
// 如果访问令牌（Access Token）尚未过期，并且调用刷新令牌（Refresh Token）时，根据常规情况，不会生成新的令牌。
//	@Summary		使用refresh token来获取新的access token
//	@Description	使用refresh token来获取新的access token
//	@Tags			Auth
//	@Accept			application/json
//	@Produce		application/json
//	@Param			refresh_token	query	string	false	"刷新Token"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]bool
//	@Router			/auth//login/refresh-token [get]
func RefreshToken(ctx *gin.Context) {
	rt := ctx.Query("refresh_token")
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(ctx, CodeInvalidToken, "缺少Auth Token")
		ctx.Abort()
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(ctx, CodeInvalidToken, "Token 格式不正确")
		ctx.Abort()
		return
	}
	accessToken, refreshToken, err := jwt.RefreshToken(parts[1], rt)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
