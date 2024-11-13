package mysql

import (
	"strings"

	"gorm.io/gorm"

	"github.com/xiaorui/reddit-async/reddit-backend/models"
)

func CreatePost(p *models.Post) (err error) {
	// 	sqlStr := `insert into post (
	// 	post_id, title, content, author_id, community_id)
	// 	values (?,?,?,?,?)
	// `
	// 	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	res := DB.Create(p)
	err = res.Error
	return
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	// sqlStr := `select
	// post_id, title, content, author_id, community_id
	// from post where post_id=?
	// `
	// res := db.Where("post_id = ?", pid).First(post)
	// res := db.First(post, "post_id = ?", pid)
	// res := db.Where("post_id = ?", pid).First(post)
	res := DB.Table("posts").Where("post_id = ?", pid).First(post)
	err = res.Error
	if post == nil {
		return nil, ErrorPostNotExist
	}

	return
}

func GetPostList(pageNum, pageSize int64) (posts []*models.Post, err error) {
	// 	sqlStr := `select
	// 	post_id, title, content, author_id, community_id, create_time
	// 	from post
	// 	ORDER BY create_time
	// 	DESC
	// 	limit ?, ?
	// `
	posts = make([]*models.Post, 0, pageSize)
	// 	err = db.Select(&posts, sqlStr, (pageNum-1)*pageSize, pageSize)
	err = DB.Model(&models.Post{}).
		Order("create_time DESC").
		Offset((int(pageNum) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&posts).Error

	return
}

func GetPostListByIDs(pids []string) (posts []*models.Post, err error) {
	// 	sqlStr := `select
	// 	post_id, title, content, author_id, community_id, create_time
	// 	from post
	// 	where post_id in (?)
	// 	ORDER BY FIND_IN_SET(post_id, ?)
	// `
	// 	query, args, err := sqlx.In(sqlStr, pids, strings.Join(pids, ","))
	// 	if err != nil {
	// 		return
	// 	}
	// 	query = db.Rebind(query) // 就是要加上?

	// err = db.Select(&posts, query, args...)
	orderArg := strings.Join(pids, ",")
	err = DB.Model(&models.Post{}).
		Where("post_id IN ? ", pids).
		Order(gorm.Expr("FIND_IN_SET(post_id, ?)", orderArg)).
		Find(&posts).Error
	return
}
