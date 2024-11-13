package logic

import (
	"github.com/xiaorui/reddit-async/reddit-backend/dao/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/snowflake"
)

func GetCommunityList() ([]*models.Community, error) {
	//其实这个业务就只有一个内容， 它不需要去判断用户是否存在， 是否登录在前面的中间件已经完成了
	return mysql.GetCommunityList()
}

func GetCommunityDetail(communityID int64) (*models.Community, error) {
	//这个函数处理的逻辑就是根据id去数据库中查询数据
	return mysql.GetCommunityDetailByID(communityID)
}

// CreateNewCommunity: 创建新的社区
func CreateNewCommunity(p *models.ParamCommunity) error {
	// 1. 查询该社区是否存在
	if err := mysql.CheckCommunityExist(p.Name); err != nil {
		return err
	}
	// 2. 创建新的社区
	// 生成uid
	uid := snowflake.GenID()

	// 3. 构造社区实例
	comm := &models.Community{
		ID:           uid,
		Name:         p.Name,
		Introduction: p.Introduction,
	}

	return mysql.InsertCommunity(comm)
}

// UpdateCommunity： 更新社区信息
func UpdateCommunity(cid string, p *models.ParamCommunity) (com *models.Community, err error) {
	// 1.  查询该community是否存在
	if com, err = mysql.GetCommunityByID(cid); err != nil {
		return nil, err
	}

	// 2. 查询更新的那个社区名是否已经存在
	if err := mysql.CheckCommunityExist(p.Name); err != nil {
		return nil, err
	}

	// 更改社区信息
	com.Name = p.Name
	com.Introduction = p.Introduction

	// 写回数据库
	return mysql.SaveCommunity(com)
}

func DeleteCommunity(cid string) error {
	return mysql.DeleteCommunity(cid)
}
