package logic

import (
	"github.com/xiaorui/reddit-async/reddit-backend/dao/redis"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
	"go.uber.org/zap"
)

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID), zap.Int64("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(userID, p.PostID, p.Direction)
}
