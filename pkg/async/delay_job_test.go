package async

import (
	"fmt"
	"testing"
	"time"
)

func TestSendAsync(t *testing.T) {
	start := time.Now().Add(1 * time.Second)
	body := map[string]interface{}{
		"type":      "profile_msg",
		"startTime": start.Format("2006-01-02 15:04:05"),
		"body": map[string]interface{}{
			"toList":  []string{"1532979219@qq.com"},
			"subject": "【测试】诚邀您更新个人信息",
			// "content": "你好的",
			"content": fmt.Sprintf(`
				<p>亲爱的 %s：</p>
				<p>我们注意到您已经加入我们社区一段时间了，感谢您的活跃参与！为了让社区成员能更好地认识您，不知您是否愿意抽出几分钟时间来丰富一下您的个人资料呢？</p>
				<p>一个精心打造的个人主页能让您：</p>
				<ul>
					<li>展示您的专业领域和特长</li>
					<li>结识更多志同道合的朋友</li>
					<li>获得更多互动和关注</li>
					<li>让您的观点和建议更具说服力</li>
				</ul>
				<p>完善个人资料非常简单，您可以：</p>
				<ol>
					<li>上传一张个性化头像</li>
					<li>添加您的兴趣爱好和专业背景</li>
					<li>写下一段独特的个性签名</li>
					<li>分享您的社交媒体链接（如果愿意的话）</li>
				</ol>
				<p><a href="[个人资料设置链接]">点击这里</a>立即开始完善您的个人主页吧！</p>
				<p>如果您在更新过程中遇到任何问题，欢迎随时联系我们的支持团队。</p>
				<p>期待在社区中看到一个全新的您！</p>
				<p>顺祝安好，</p>
				<p>社区运营团队</p>
			`, "测试"),
			"user_id": "572723818020212736",
		},
		"ttl":   100,
		"retry": 0,
	}
	err := SendAsyncTask(body)
	if err != nil {
		t.Fatal("执行错误")
	}
}
