package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/settings"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var db *sqlx.DB
var (
	DB    *gorm.DB
	SQLDB *sql.DB
)

func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	fmt.Println(cfg.Port)
	// 也可以使用MustConnect连接不成功就panic
	// db, err = sqlx.Connect("mysql", dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// fmt.Printf("connect DB failed, err:%v\n", err)
		zap.L().Error("connect DB failed.", zap.Error(err), zap.Any("端口是：", cfg.Port))
		return
	}

	SQLDB, err = DB.DB()
	SQLDB.SetMaxOpenConns(cfg.MaxOpenConns) // 设置最大的连接数量
	SQLDB.SetMaxIdleConns(cfg.MaxIdleConns) // 设置最大的空闲连接的数量， 为了避免空闲连接占用资源

	// TODO:这里写数据库迁移的操作，后面进行更新
	DB.AutoMigrate(&models.User{}, &models.Community{}, &models.Post{}, &models.Comment{}) // 会默认使用复数形式
	return
}

func Close() {
	_ = SQLDB.Close()
}
