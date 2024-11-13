package mysql

import (
	"testing"

	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/settings"
)

func init() {
	dbCfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		Port:         13306,
		User:         "root",
		Password:     "root1234",
		DBName:       "bluebell",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	p := &models.Post{
		ID:          10,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(p)
	if err != nil {
		t.Fatalf("create insert record into mysql failed, err : %v\n", err)
	}

	t.Logf("create insert record into mysql success")
}
