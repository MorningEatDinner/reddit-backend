package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/mail"
	"github.com/xiaorui/reddit-async/reddit-backend/settings"
)

func Consumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"email",  // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,    // queue name
		"welcome", // routing key
		"email",   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	forever := make(chan bool)
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	go func() {
		for d := range msgs {
			body := d.Body
			p := &models.ParamSignUpUsingEmail{}
			json.Unmarshal(body, p)
			_ = mail.NewMailer().Send(
				mail.Email{
					From: mail.From{
						settings.Conf.EmailConfig.FromConfig.Address, // 发件人的地址
						settings.Conf.EmailConfig.FromConfig.Name,    // 名称
					},
					To:      []string{p.Email},                     // 收件人地址
					Subject: fmt.Sprintf("欢迎加入Reddit, %s", p.Name), // 主题
					HTML: []byte(fmt.Sprintf(`
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
					`, p.Name)), // 内容
				},
			)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
