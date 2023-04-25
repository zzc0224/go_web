package redis

import (
	"bluebell_backend/dao/mysql"
	"math"
	"time"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

const (
	OneYearInSeconds         = 365 * 7 * 24 * 3600
	VoteScore        float64 = 432 //  86400/200=432
	PostPerAge               = 20
)

var db *sqlx.DB

/*
投票算法：http://www.ruanyifeng.com/blog/2012/03/ranking_algorithm_reddit.html



*/

/*
	PostVote 为帖子投票

投票分为四种情况：1.投赞成票 2.投反对票 3.取消投票 4.反转投票

记录文章参与投票的人
更新文章分数：赞成票要加分；反对票减分

v=1时，有两种情况

	1.之前没投过票，现在要投赞成票
	2.之前投过反对票，现在要改为赞成票

v=0时，有两种情况

	1.之前投过赞成票，现在要取消
	2.之前投过反对票，现在要取消

v=-1时，有两种情况

	1.之前没投过票，现在要投反对票
	2.之前投过赞成票，现在要改为反对票
*/
func PostVote(postID, userID string, v float64) (err error) {
	// 1. 取帖子发布时间
	postTime := client.ZScore(KeyPostTimeZSet, postID).Val()
	if float64(time.Now().Unix())-postTime > OneYearInSeconds {
		// 不允许投票了
		return ErrorVoteTimeExpire
	}
	// 判断是否已经投过票
	key := KeyPostVotedZSetPrefix + postID
	ov := client.ZScore(key, userID).Val() // 获取当前分数

	diffAbs := math.Abs(ov - v)
	pipeline := client.TxPipeline()
	pipeline.ZAdd(key, redis.Z{ // 记录已投票
		Score:  v,
		Member: userID,
	})
	pipeline.ZIncrBy(KeyPostScoreZSet, VoteScore*diffAbs*v, postID) // 更新分数
	//switch math.Abs(ov) - math.Abs(v) {
	//case 1:
	//	// 取消投票 ov=1/-1 v=0
	//	// 投票数-1
	//	pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", -1)
	//case 0:
	//	// 反转投票 ov=-1/1 v=1/-1
	//	// 投票数不用更新
	//case -1:
	//	// 新增投票 ov=0 v=1/-1
	//	// 投票数+1
	//	pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", 1)
	//default:
	//	// 已经投过票了
	//	return ErrorVoted
	//}
	switch ov - v {
	case -2:
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", 2)
	case -1:
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", 1)
	case 0:
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", -int64(v))
		pipeline.ZAddXX(key, redis.Z{ // 将记录中的投票取消
			Score:  0,
			Member: userID,
		})
		pipeline.ZIncrBy(KeyPostScoreZSet, -VoteScore*v, postID) // 更新分数
	case 1:
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", -1)
	case 2:
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", -2)
	}
	_, err = pipeline.Exec()
	zap.L().Debug("PostVote", zap.String("postID", postID), zap.String("userID", userID), zap.Float64("direction", v))
	return
}

// CreatePost 使用hash存储帖子信息
func CreatePost(postID, userID, title, summary, communityName string) (err error) {
	now := float64(time.Now().Unix())
	votedKey := KeyPostVotedZSetPrefix + postID
	communityKey := KeyCommunityPostSetPrefix + communityName
	postInfo := map[string]interface{}{
		"title":     title,
		"summary":   summary,
		"post:id":   postID,
		"user:id":   userID,
		"time":      now,
		"votes":     1,
		"comments":  0,
		"community": communityName,
	}

	// 事务操作
	pipeline := client.TxPipeline()
	pipeline.ZAdd(votedKey, redis.Z{ // 作者默认投赞成票
		Score:  1,
		Member: userID,
	})
	pipeline.Expire(votedKey, time.Second*OneYearInSeconds) // 一年时间

	pipeline.HMSet(KeyPostInfoHashPrefix+postID, postInfo)
	pipeline.ZAdd(KeyPostScoreZSet, redis.Z{ // 添加到分数的ZSet
		Score:  now + VoteScore,
		Member: postID,
	})
	pipeline.ZAdd(KeyPostTimeZSet, redis.Z{ // 添加到时间的ZSet
		Score:  now,
		Member: postID,
	})
	pipeline.SAdd(communityKey, postID) // 添加到对应版块
	_, err = pipeline.Exec()
	return
}

// GetPost 从key中分页取出帖子(一页PostPerAge 20 个帖子)
func GetPost(order string, page int64) []map[string]string {
	key := KeyPostScoreZSet
	if order == "time" {
		key = KeyPostTimeZSet
	}
	start := (page - 1) * PostPerAge
	end := start + PostPerAge - 1
	ids := client.ZRevRange(key, start, end).Val()
	postList := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		postData := client.HGetAll(KeyPostInfoHashPrefix + id).Val()
		postData["id"] = id
		postList = append(postList, postData)
	}
	return postList
}

func GetReCommendList(keys []string) []map[string]string {
	recommendList := make([]map[string]string, 0)
	for _, key := range keys {
		recommend := client.HGetAll(KeyPostInfoHashPrefix + key).Val()
		recommend["id"] = key
		recommendList = append(recommendList, recommend)
	}
	return recommendList
}

func GetVote() map[string]map[string]float64 {
	recommendMap := make(map[string]map[string]float64)
	var PostIdList []string
	postKey := KeyPostVotedZSetPrefix + "*"
	result, _ := client.Keys(postKey).Result()
	for _, val := range result {
		PostIdList = append(PostIdList, val)
	}
	userIdList := mysql.GetAllUser()
	for _, userId := range userIdList {
		recommendMap[userId] = make(map[string]float64)
	}
	for _, s := range PostIdList {
		for _, userId := range userIdList {
			recommendMap[userId][s] = 0
		}
		//println("vote value")
	}
	for _, s := range PostIdList {
		//println(i, s)
		strings, _ := client.ZRangeWithScores(s, 0, -1).Result()

		for _, z := range strings {
			//println(z.Score, z.Member)
			//fmt.Printf("%v %v\n", z.Score, z.Member)
			userId := z.Member.(string)
			recommendMap[userId][s] = z.Score
		}
		//println("vote value")
	}
	return recommendMap
}

// GetCommunityPost 分社区根据发帖时间或者分数取出分页的帖子
func GetCommunityPost(communityName, orderKey string, page int64) []map[string]string {
	key := orderKey + communityName // 创建缓存键

	if client.Exists(key).Val() < 1 {
		client.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, KeyCommunityPostSetPrefix+communityName, KeyPostInfoHashPrefix+orderKey)
		client.Expire(key, 60*time.Second)
	}
	start := (page - 1) * PostPerAge
	end := start + PostPerAge - 1
	ids := client.ZRevRange(key, start, end).Val()
	postList := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		postData := client.HGetAll(KeyPostInfoHashPrefix + id).Val()
		postData["id"] = id
		postList = append(postList, postData)
	}
	return postList

	//return GetPost(key, page)
}

// Reddit Hot rank algorithms
// from https://github.com/reddit-archive/reddit/blob/master/r2/r2/lib/db/_sorts.pyx
func Hot(ups, downs int, date time.Time) float64 {
	s := float64(ups - downs)
	order := math.Log10(math.Max(math.Abs(s), 1))
	var sign float64
	if s > 0 {
		sign = 1
	} else if s == 0 {
		sign = 0
	} else {
		sign = -1
	}
	seconds := float64(date.Second() - 1577808000)
	return math.Round(sign*order + seconds/43200)
}
