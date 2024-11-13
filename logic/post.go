package logic

import (
	"errors"

	"github.com/xiaorui/reddit-async/reddit-backend/dao/mysql"
	"github.com/xiaorui/reddit-async/reddit-backend/dao/redis"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/snowflake"
	"go.uber.org/zap"
)

var (
	ErrorNotPerm = errors.New("无权限操作")
)

func CreatePost(p *models.Post) (err error) {
	//1. 生成post id
	p.ID = snowflake.GenID()

	//2. 将数据保存到数据库, 这里还需要再redis中加入post的记录， 当前post创建的时间
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}

	err = redis.CreatePost(p.ID, p.CommunityID, p.AuthorID)
	if err != nil {
		return
	}
	return
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	//就是从mysql中去获取数据
	//不只需要post的信息， 还需要community的信息，还需要author的信息
	//1. 先获取post， 才能获取author， 才能获取community
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostByID mysql.GetPostByID failed.", zap.Error(err))
		return nil, err
	}
	//2. 获取authorname
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("GetPostByID mysql.GetUserByID failed.", zap.Error(err))
		return nil, err
	}

	//3. 获取community信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("GetPostByID mysql.GetCommunityDetailByID failed.", zap.Error(err))
		return nil, err
	}

	data = &models.ApiPostDetail{
		user.Username,
		post,
		community,
	}

	return
}

func GetPostList(pageNum, pageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, pageSize)
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//2. 获取authorname
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("GetPostByID mysql.GetUserByID failed.", zap.Error(err))
			continue
		}

		//3. 获取community信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetPostByID mysql.GetCommunityDetailByID failed.", zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			user.Username,
			post,
			community,
		}
		data = append(data, postDetail)
	}

	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail2, err error) {
	//1. 先去redis查询得到post id的列表。 显示的先后依据是从这里来看的。
	pidList, err := redis.GetPostIDListByOrder(p)
	if err != nil {
		return nil, err
	}
	//如果pidList是空的
	if len(pidList) == 0 {
		zap.L().Warn("GetPostList2 len(pidList) == 0")
		return
	}

	//从redis中去获取这些pids的票数
	votes, err := redis.GetVotesByPostIDS(pidList)
	if err != nil {
		return nil, err
	}

	//2. 根据列表从数据库中得到post的详细信息
	posts, err := mysql.GetPostListByIDs(pidList)
	if err != nil {
		return nil, err
	}

	//就是说希望在这里传入用户对于每个帖子的投票情况， 应该在结构体中加入一个结构信息， 即投票
	//3. 根据获取到的post的详细信息去mysql中查询community和user的详细信息
	data = make([]*models.ApiPostDetail2, 0, len(posts))
	for idx, post := range posts {
		//2. 获取authorname
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("GetPostByID mysql.GetUserByID failed.", zap.Error(err))
			continue
		}

		//3. 获取community信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetPostByID mysql.GetCommunityDetailByID failed.", zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail2{
			user.Username,
			votes[idx],
			post,
			community,
		}
		data = append(data, postDetail)
	}

	return
}

// 这个函数的主要目的就是加上communityid， 也就是说获取pid的这里的方式需要有community的参与
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail2, err error) {
	//1. 先去redis查询得到post id的列表
	//pidList, err := redis.GetPostIDListByOrder(p)
	pidList, err := redis.GetCommunityPostIDListByOrder(p)
	if err != nil {
		return nil, err
	}
	//如果pidList是空的
	if len(pidList) == 0 {
		zap.L().Warn("GetPostList2 len(pidList) == 0")
		return
	}

	//从redis中去获取这些pids的票数
	votes, err := redis.GetVotesByPostIDS(pidList)
	if err != nil {
		return nil, err
	}
	//2. 根据列表从数据库中得到post的详细信息
	posts, err := mysql.GetPostListByIDs(pidList)
	if err != nil {
		return nil, err
	}

	//就是说希望在这里传入用户对于每个帖子的投票情况， 应该在结构体中加入一个结构信息， 即投票
	//3. 根据获取到的post的详细信息去mysql中查询community和user的详细信息
	data = make([]*models.ApiPostDetail2, 0, len(posts))
	for idx, post := range posts {
		//2. 获取authorname
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("GetPostByID mysql.GetUserByID failed.", zap.Error(err))
			continue
		}

		//3. 获取community信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetPostByID mysql.GetCommunityDetailByID failed.", zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail2{
			user.Username,
			votes[idx],
			post,
			community,
		}
		data = append(data, postDetail)
	}

	return
}

func GetPostList0(p *models.ParamPostList) (data []*models.ApiPostDetail2, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostList0 failed.", zap.Error(err))
		return nil, err
	}
	return
}

// DeletePost: 删除post
func DeletePost(postID, userID int64) error {
	// 1. 查询post
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		return err
	}

	// 2. 执行删除操作
	if err := mysql.DeletePost(post, userID); err != nil {
		return err
	}

	return nil
}
