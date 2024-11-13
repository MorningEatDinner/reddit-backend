package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/async"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/rabbitmq"
)

func TestAsync(ctx *gin.Context) {
	// 1. 发送异步邮件请求-- 使用这个异步系统
	now := time.Now()
	start := now.Add(1 * time.Second)
	body := map[string]interface{}{
		"type":      "mail_msg",
		"startTime": start.Format("2006-01-02 15:04:05"),
		"body": map[string]interface{}{
			"toList":  []string{"12984532492@qq.com"},
			"subject": "欢迎加入Reddit社区",
			"content": fmt.Sprintf(`
					<p>亲爱的 %s，您好！</p>
					<p>感谢您注册并加入 Reddit！我们非常高兴能有机会与您一同分享这里的丰富内容与有趣讨论。无论您是来学习新知识、分享经验，还是寻找志同道合的朋友，我们相信您一定能在这里找到属于自己的精彩世界。</p>
					<p>为了帮助您更好地融入社区，以下是一些小提示：</p>
					<ol>
						<li>完善个人资料：登录后，您可以前往个人资料页面，添加头像和个性签名，让大家更好地认识您。</li>
						<li>阅读社区规则：为了营造一个友好、和谐的讨论环境，请花一点时间阅读我们的<a href="[社区规则链接]">社区规则链接</a>。</li>
						<li>参与讨论：不要害羞，随时可以加入我们已有的热门话题，或者发起您感兴趣的讨论。</li>
						<li>如果您在使用过程中有任何问题，随时可以联系我们的客服团队或管理员，我们将竭诚为您提供帮助。</li>
					</ol>
					<p>再次感谢您的加入，期待与您在社区中互动！</p>
					<p>祝好，</p>
				`, "test"),
		},
		"ttl":   100,
		"retry": 1,
	}
	if err := async.SendAsyncTask(body); err != nil {
		ResponseError(ctx, CodePhoneCodeSendError)
	}
	ResponseSuccess(ctx, nil)
}

func TestMq(ctx *gin.Context) {
	//  1. 进行参数验证
	p := new(models.ParamSignUpUsingEmail)
	if err := ctx.ShouldBindJSON(p); err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	// 2. 进行业务处理, 创建新的用户
	if err := rabbitmq.PublishEmailTask(p); err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(ctx, nil)
}
