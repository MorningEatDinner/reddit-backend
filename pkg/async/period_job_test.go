package async

import (
	"testing"
)

func TestSendPeriod(t *testing.T) {
	cronSpec := "0/5 * * * *"
	body := map[string]interface{}{
		"type":     "mail_msg",
		"cronSpec": cronSpec,
		"body": map[string]interface{}{
			"toList":  []string{"1532979219@qq.com"},
			"subject": "这是对于周期任务的测试",
			"content": "你好的",
		},
		"ttl":   100,
		"retry": 0,
	}
	err := SendPeriodTask(body)
	if err != nil {
		t.Fatal("执行错误")
	}
}
