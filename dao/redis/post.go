package redis

import (
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xiaorui/reddit-async/reddit-backend/models"
)

func CreatePost(pid, communityID, userID int64) error {
	// 在redis中记录四条数据：1。 帖子的分数， 时 用户的给帖子的透片社区有纳西额帖子
	// Redis事务虽然不满足acid属性， 但是能够满足部分的原子性， 要么全部执行， 要么全部不执行
	//使用事务操作， redis事务
	pipeline := RDB.Client.TxPipeline()
	// 作者默认投票投赞成票, 这里只是记录投票的方向
	votedKey := getRedisKey(KeyPostVotedZSetPF) + strconv.Itoa(int(pid))
	pipeline.ZAdd(RDB.Context, votedKey, redis.Z{
		Score:  1,
		Member: userID,
	})

	//在redis中加入一个创建的post的记录
	pipeline.ZAdd(RDB.Context, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid, //需要确认这里是否需要是string类型的？先暂时使用int
	})
	// 在创建post的时候， 也加入了对于分数的初始化
	pipeline.ZAdd(RDB.Context, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid, //需要确认这里是否需要是string类型的？先暂时使用int
	})

	communityKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	// 我们也可以在创建comment的地方采用相同的策略，创建一个postid的key， 然后使用Redis来存放有哪些commentid， 这样查询的时候就无需逐条查询mysql了
	pipeline.SAdd(RDB.Context, communityKey, pid) // 加入member， 但是不需要score， 就是给community下面添加数据， 这些数据是使用Set来保存的

	_, err := pipeline.Exec(RDB.Context)
	return err
}

func getIDSFromKey(key string, page, size int64) ([]string, error) {
	//得到按照某种方式进行排序
	start := (page - 1) * size
	end := page * size
	// 根据分数或者时间从高到低去获取部分的pid
	return RDB.Client.ZRevRange(RDB.Context, key, start, end).Result() // 从高到低
}

func GetPostIDListByOrder(p *models.ParamPostList) ([]string, error) {
	key := getRedisKey(KeyPostTimeZSet) // 根据时间拿到数据
	if p.Order == models.OrderScore {   // 根据投票分数拿到数据
		key = getRedisKey(KeyPostScoreZSet)
	}

	return getIDSFromKey(key, p.Page, p.Size)
}

// GetVotesByPostIDS: 得到帖子的投票数
func GetVotesByPostIDS(pidList []string) ([]int64, error) {
	pipeline := RDB.Client.Pipeline()
	for _, id := range pidList {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(RDB.Context, key, "1", "1") // 统计区间在1都1之间的数据是多少。 统计给正票的数量
	}
	cmders, err := pipeline.Exec(RDB.Context)
	if err != nil {
		return nil, err
	}
	data := make([]int64, 0, len(pidList))
	for _, cmder := range cmders {
		value := cmder.(*redis.IntCmd).Val()
		data = append(data, value)
	}
	return data, nil
}

func GetCommunityPostIDListByOrder(p *models.ParamPostList) ([]string, error) {
	// 再多加上一个key， 如果一段时间内重复查询会更快， 也就是加上一个对之前查询结果的缓存
	communityKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	key := orderKey + strconv.Itoa(int(p.CommunityID)) // 查找某个社区的post，按照order排序
	// 如果不存在， 也就是如果缓存里面没有， 那么就需要查询了
	if RDB.Client.Exists(RDB.Context, key).Val() < 1 { //
		// 需要计算
		pipeline := RDB.Client.Pipeline()
		pipeline.ZInterStore(RDB.Context, key, &redis.ZStore{
			Aggregate: "MAX",                            // 这里的意思是相同元素的聚合方式
			Keys:      []string{communityKey, orderKey}, // 计算两个有序集合的交集
		}) // 注意， 值最终是保存到一个zset中的
		pipeline.Expire(RDB.Context, key, time.Second*60) // 只有60秒的生存时间， 因为实时性要求高吗
		_, err := pipeline.Exec(RDB.Context)
		if err != nil {
			return nil, err
		}
	}

	// 上面结束之后就得到key中的值就是获取了key的对应的值， 就是说所有的id
	// 这里就是按照某种分页的依据来实现数据的获取
	return getIDSFromKey(key, p.Page, p.Size)
}
