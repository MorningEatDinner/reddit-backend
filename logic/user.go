package logic

import (
	"errors"
	"fmt"
	"time"

	"github.com/xiaorui/reddit-async/reddit-backend/dao/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/dao/redis"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/async"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/helpers"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/jwt"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/mail"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/rabbitmq"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/sms"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/snowflake"
	"github.com/xiaorui/reddit-async/reddit-backend/settings"
	"go.uber.org/zap"
)

// 存放业务逻辑的代码
func SignUp(p *models.ParamSignUp) (err error) {
	//1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//2. 生成uid
	userID := snowflake.GenID()
	//3. 密码加密

	//构造数据实例
	user := &models.User{
		ID:       userID,
		Username: p.Username,
		Password: p.Password,
	}

	//4. 保存进数据库

	err = mysql.InsertUser(user)

	//这里还可以有很多其他的数据操作， 比如对于redis进行操作

	return
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(user); err != nil {
		return nil, err
	}
	//如果登录成功
	//生成JWT
	//return jwt.GenToken(user.UserID, user.Username)
	accessToken, _, err := jwt.GenToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = accessToken
	return
}

// LoginUsingPhoneWithCode: 使用手机+验证码的形式进行登陆
func LoginUsingPhoneWithCode(p *models.ParamLoginUsingPhoneWithCode) (*models.User, error) {
	user := &models.User{
		Phone: p.Phone,
	}
	// 进行登陆操作
	if err := mysql.LoginUsingPhoneWithCode(user); err != nil {
		return nil, err
	}

	// 如果登录成功：即确实有这个用户存在
	return user, nil
}

// LoginUsingEmail: 使用邮箱+密码的方式进行登陆
func LoginUsingEmail(p *models.ParamLoginUsingEmail) (*models.User, error) {
	user := &models.User{
		Email:    p.Email,
		Password: p.Password,
	}
	if err := mysql.LoginUsingEmail(user); err != nil {
		zap.L().Error(" mysql.LoginUsingEmail failed", zap.Error(err))
		return nil, err
	}

	return user, nil
}

// IsPhoneExist：返回输入手机号码是否存在数据表中
func IsPhoneExist(phone string) (bool, error) {
	exist, err := mysql.IsPhoneExist(phone)
	if err == mysql.ErrorPhoneExist {
		return true, nil
	}
	return exist, err
}

// IsEmailExist：返回输入邮箱是否存在数据表中
func IsEmailExist(phone string) (bool, error) {
	return mysql.IsEmailExist(phone)
}

// SendPhoneCode: 发送短信验证码
func SendPhoneCode(phone string) error {
	// 1. 生成验证码
	code := helpers.GenerateRandomCode()

	// 2. 将验证码保存到redis中
	if err := redis.SetVerifyCode(phone, code); err != nil {
		zap.L().Error("SetVerifyCode failed...", zap.Error(err))
		return err
	}
	// 3. 发送短信给手机
	if ok := sms.NewSms().Send(phone, code); !ok {
		return errors.New("短信发送失败")
	}
	return nil
}

// SendEmailCode: 发送邮箱验证码
func SendEmailCode(email string) error {
	// 1. 生成验证码
	code := helpers.GenerateRandomCode()

	// 2. 将验证码存放到redis中
	if err := redis.SetVerifyCode(email, code); err != nil {
		zap.L().Error("SetVerifyCode failed...", zap.Error(err))
		return err
	}
	// 3. 发送验证码给邮箱;
	ok := mail.NewMailer().Send(
		mail.Email{
			From: mail.From{
				settings.Conf.EmailConfig.FromConfig.Address, // 发件人的地址
				settings.Conf.EmailConfig.FromConfig.Name,    // 名称
			},
			To:      []string{email},                                         // 收件人地址
			Subject: "email 验证码",                                             // 主题
			HTML:    []byte(fmt.Sprintf("<h1>您的 Email 验证码是 %v </h1>", code)), // 内容
		},
	)
	if !ok {
		return errors.New("邮箱验证码发送失败")
	}
	return nil
}

// SignupUsingPhone：处理手机注册登陆逻辑
func SignupUsingPhone(p *models.ParamSignupUsingPhone) (err error) {
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Name); err != nil {
		return err
	}

	// 判断手机号码是否存在
	if _, err = mysql.IsPhoneExist(p.Phone); err != nil {
		return err
	}
	// 2. 生成uid
	userID := snowflake.GenID()

	// 3. 构造用户实例
	user := &models.User{
		ID:       userID,
		Username: p.Name,
		Password: p.Password,
		Phone:    p.Phone,
	}

	// 4. 保存到数据库
	err = mysql.InsertUser(user)

	return
}

// SignUpUsingEmail: 进行使用邮箱进行注册的业务
func SignUpUsingEmail(p *models.ParamSignUpUsingEmail) (err error) {
	// 1. 验证用户名是否存在
	if err = mysql.CheckUserExist(p.Name); err != nil {
		return err
	}
	// 2. 验证邮箱是否已经注册
	if _, err = mysql.IsEmailExist(p.Email); err != nil {
		return err
	}
	// 3. 生成uid
	userID := snowflake.GenID()

	// 4. 构造用户实例
	_user := models.User{
		ID:       userID,
		Username: p.Name,
		Email:    p.Email,
		Password: p.Password,
	}

	// TODO: 这里用户注册成功之后， 发送一份邮件给用户邮箱， 说明欢迎加入论坛
	rabbitmq.PublishEmailTask(p)

	// 注册成功之后投递一个任务，向邮箱中发送欢迎加入社区的邮件
	now := time.Now()
	start := now.Add(1 * time.Second)
	body := map[string]interface{}{
		"type":      "mail_msg",
		"startTime": start.Format("2006-01-02 15:04:05"),
		"body": map[string]interface{}{
			"toList":  []string{p.Email},
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
				`, p.Name),
		},
		"ttl":   100,
		"retry": 0,
	}
	if err = async.SendAsyncTask(body); err != nil {
		return err
	}

	if err = redis.SetProfileStatus(fmt.Sprint(userID), false); err != nil {
		return err
	}
	// 发起一个定时任务， 如果24小时之后查看这个键还是没有修改， 则向邮箱发送一个邮件，邀请修改个人信息
	start = now.Add(24 * time.Second)
	body = map[string]interface{}{
		"type":      "profile_msg",
		"startTime": start.Format("2006-01-02 15:04:05"),
		"body": map[string]interface{}{
			"toList":  []string{p.Email},
			"subject": "诚邀您更新个人信息",
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
			`, p.Name),
			"user_id": fmt.Sprintf("%v", userID),
		},
		"ttl":   100,
		"retry": 0,
	}
	if err = async.SendAsyncTask(body); err != nil {
		return err
	}
	// 5. 保存到数据库中
	if err = mysql.InsertUser(&_user); err != nil {
		return err
	}
	return nil
}

// UpdateProfile: 修改用户名+简介+城市
func UpdateProfile(p *models.ParamUpdateProfile, userID int64) (user *models.User, err error) {
	// 1. 查询当前用户
	if user, err = mysql.GetUserByID(userID); err != nil {
		return nil, err
	}

	// 2. 如果要更改的用户名不是是当前用户名， 则查看用户名是否存在
	if user.Username != p.Name {
		if err = mysql.CheckUserExist(p.Name); err != nil {
			return nil, err
		}
	}

	// 3. 设置用户信
	user.Username = p.Name
	user.City = p.City
	user.Introduction = p.Introduction
	// 4. 写回数据库
	user, err = mysql.SaveUser(user)
	if err != nil {
		return nil, err
	}
	// 更新用户的 更新状态
	err = redis.SetProfileStatus(fmt.Sprintf("%v", userID), true)
	if err != nil {
		return nil, err
	}
	return
}

// UpdateEmail: 修改用户邮箱
func UpdateEmail(p *models.ParamUpdateEmail, userID int64) (user *models.User, err error) {
	// 1. 查询当前用户
	if user, err = mysql.GetUserByID(userID); err != nil {
		return nil, err
	}

	// 2. 查看当下邮箱是否存在
	if _, err = mysql.IsEmailExist(p.Email); err != nil {
		return nil, err
	}

	// 3. 设置用户信
	user.Email = p.Email
	// 4. 写回数据库
	return mysql.SaveUser(user)
}

func UpdatePhone(p *models.ParamUpdatePhone, userID int64) (user *models.User, err error) {
	// 1. 查询当前用户
	if user, err = mysql.GetUserByID(userID); err != nil {
		return nil, err
	}
	// 2. 查看当前号码是否存在
	if _, err = mysql.IsPhoneExist(p.Phone); err != nil {
		return nil, err
	}

	// 3. 设置新的号码
	user.Phone = p.Phone

	// 4. 进行写回
	return mysql.SaveUser(user)
}

// UpdatePassword： 更改当前用户的密码
func UpdatePassword(p *models.ParamUpdatePassword, userID int64) (err error) {
	return mysql.UpdatePassword(p.Password, p.NewPassword, userID)
}

func GetEmailList() (dataList []string, err error) {
	return mysql.GetEmailList()
}
