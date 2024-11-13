package mysql

import (
	"crypto/md5"
	"encoding/hex"

	"gorm.io/gorm"

	"github.com/xiaorui/reddit-async/reddit-backend/models"
)

const secret = "xiaorui.com"

// CheckUserExist: 检查制定用户名是否存在
func CheckUserExist(username string) (err error) {
	var count int64
	// sqlStr := `select count(user_id) from user where username = ?`
	res := DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if err = res.Error; err != nil {
		return
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser: 向数据库中添加一条新的用户数据
func InsertUser(user *models.User) (err error) {
	//在插入数据前， 需要对密码进行加密处理
	user.Password = encryptPassword(user.Password)
	// 将数据实例插入数据表中
	// sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	// _, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	res := DB.Create(user)
	err = res.Error
	return
}

// 使用md5加密密码现在被认为是不安全的了， 因为可以通过如暴力破解的方式来破解
func encryptPassword(oPasssword string) string {
	h := md5.New()                                       // 创建hash.Hash对象
	h.Write([]byte(secret))                              //  写入颜值
	return hex.EncodeToString(h.Sum([]byte(oPasssword))) // 对于数据进行加密
}

func Login(user *models.User) (err error) {
	//这里进行用户登陆， 也就是根据用户名查询用户数据， 之后验证密码是否相等，完成
	// sqlStr := `select user_id, username, password from user where username=?`
	oPassword := user.Password // 记录原始密码， 后面从数据库中返回的密码会叠加在这个数据上
	// err = db.Get(user, sqlStr, user.Username)
	err = DB.Model(&models.User{}).Where("username = ?", user.Username).First(user).Error
	if err == gorm.ErrRecordNotFound {
		//var ErrNoRows = errors.New("sql: no rows in result set")
		//想要特别写出来这个错误，否则则原始错误信息返回
		return ErrorUserNotExist
	}
	if err != nil {
		return
	}
	if encryptPassword(oPassword) != user.Password {
		return ErrorPasswordInvalid
	}
	return
}

// LoginUsingPhoneWithCode： 使用手机+验证码登陆
func LoginUsingPhoneWithCode(user *models.User) error {
	// 这里因为已经验证过验证码了， 所以只需要知道用户是否存在, 即该手机号码是否注册
	// 得拿到整个user，并且写回
	err := DB.Model(&models.User{}).Where("phone = ?", user.Phone).First(user).Error
	if err == gorm.ErrRecordNotFound {
		// 如果手机号码不存在
		return ErrorPhoneNotExist
	}
	return err
}

// LoginUsingEmail: 使用email来查询用户是否存在， 并且确认密码是否一致
func LoginUsingEmail(user *models.User) error {
	// 1. 查询用户
	oPassword := user.Password
	err := DB.Model(&models.User{}).Where("email = ?", user.Email).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return ErrorEmailNotExist
	}
	if err != nil {
		return err
	}
	if encryptPassword(oPassword) != user.Password {
		return ErrorPasswordInvalid
	}
	return err
}

func GetUserByID(id int64) (user *models.User, err error) {
	user = new(models.User)
	// sqlStr := `select user_id, username from user where user_id=?`

	// err = db.Get(user, sqlStr, id)
	err = DB.Model(&models.User{}).Where("user_id = ?", id).First(user).Error

	return
}

// IsPhoneExist： 验证手机号是否已经注册了
func IsPhoneExist(phone string) (bool, error) {
	var count int64
	err := DB.Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsEmailExist: 验证邮箱是否注册了
func IsEmailExist(email string) (bool, error) {
	var count int64
	err := DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SaveUser: 保存修改后的用户信息到数据库中
func SaveUser(user *models.User) (*models.User, error) {
	res := DB.Save(user)
	if res.RowsAffected == 0 {
		return nil, ErrorSaveUser
	}
	return user, nil
}

// CheckPasswordValid：验证用户密码是否有效
func UpdatePassword(password, NewPassword string, userID int64) error {
	var user models.User
	err := DB.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return ErrorUserNotExist
	}
	if encryptPassword(password) != user.Password {
		return ErrorPasswordInvalid
	}
	user.Password = encryptPassword(NewPassword)
	_, err = SaveUser(&user)

	return err
}

// DeletePost: 删除Post
func DeletePost(post *models.Post, userID int64) error {
	if post.AuthorID != userID {
		return ErrorNotPermission
	}

	return DB.Delete(post).Error
}

func GetEmailList() ([]string, error) {
	var emailList []string
	// 查询所有用户的email
	err := DB.Model(&models.User{}).
		Select("email").
		Pluck("email", &emailList).
		Error

	return emailList, err
}
