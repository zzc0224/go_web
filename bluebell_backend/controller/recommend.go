package controller

import (
	"bluebell_backend/dao/redis"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReCommend(c *gin.Context) {
	recommendMap := redis.GetVote()
	UserID, _ := getCurrentUserID(c)
	CurrentUserID := strconv.FormatUint(UserID, 10)
	CurrentVoteMap := make(map[string]float64)

	sim := make(map[string]map[string]float64)
	sim[CurrentUserID] = make(map[string]float64)
	averageArr := make(map[string]float64)
	var CurrentAve float64
	type VoteStruct struct {
		postId string
		vote   float64
	}
	var VoteStructArr []VoteStruct
	type OtherVoteStruct struct {
		postId string
		vote   float64
	}
	var OtherVoteStructArr []OtherVoteStruct
	type reCommendIdListStruct struct {
		postKey string
		simVote float64
	}
	var reCommendIdList []reCommendIdListStruct
	var reCommendKeyList []string
	//println("test--------------------")
	//for user, list := range recommendMap {
	//	for post, vote := range list {
	//		println(user, post, vote)
	//	}
	//}
	//println("test--------------------")

	for CurrentPostId, CurrentVote := range recommendMap[CurrentUserID] {
		//println(CurrentUserID, CurrentPostId, CurrentVote)
		CurrentVoteMap[CurrentPostId] = CurrentVote
		CurrentAve += CurrentVote
	}
	CurrentAve = CurrentAve / float64(len(CurrentVoteMap))
	for CurrentPostId, CurrentVote := range CurrentVoteMap {
		if CurrentVote == 0 {
			for s, f := range CurrentVoteMap { //创建临时VoteArr
				if s != CurrentPostId {
					VoteStructArr = append(VoteStructArr, VoteStruct{s, f})
				}
			}
			for userId, List := range recommendMap {
				var x, y, z, num float64
				if userId != CurrentUserID {
					for postId, vote := range List { //创建临时OtherVoteArr
						if postId != CurrentPostId {
							OtherVoteStructArr = append(OtherVoteStructArr, OtherVoteStruct{postId, vote})
						}
					}
					length := len(OtherVoteStructArr)
					for i := 0; i < length; i++ { //由于golang的forrange不固定顺序，因此需要排序计算
						for j := 0; j < length; j++ {
							if VoteStructArr[i].postId == OtherVoteStructArr[j].postId {
								x += VoteStructArr[i].vote * OtherVoteStructArr[j].vote
								y += VoteStructArr[i].vote * VoteStructArr[i].vote
								z += OtherVoteStructArr[i].vote * OtherVoteStructArr[i].vote
								num += OtherVoteStructArr[i].vote
								//println(VoteStructArr[i].vote, OtherVoteStructArr[j].vote)
							}
						}

					}
					if y == 0 || z == 0 {
						sim[CurrentUserID][userId] = 0
					} else {
						sim[CurrentUserID][userId] = x / (math.Sqrt(y) * math.Sqrt(z)) //计算当前用户与其他用户的余弦相似度
					}
					averageArr[userId] = num / float64(length)
					OtherVoteStructArr = nil
					//CurrentVoteMap[CurrentPostId] = sim[CurrentUserID][userId] * recommendMap[userId][CurrentPostId] //根据余弦相似度修改vote,需完善
				}
			}
			VoteStructArr = nil
			//fmt.Printf("%f\n", CurrentVoteMap[CurrentPostId])
			type simStruct struct {
				user   string
				simNum float64
			}
			var simStructArr []simStruct
			for userId, simNum := range sim[CurrentUserID] {
				simStructArr = append(simStructArr, simStruct{userId, simNum})
			}
			sort.Slice(simStructArr, func(i, j int) bool {
				return simStructArr[i].simNum > simStructArr[j].simNum
			})
			if len(CurrentVoteMap) <= 4 {
				CurrentVoteMap[CurrentPostId] = CurrentAve + (recommendMap[simStructArr[0].user][CurrentPostId] - averageArr[simStructArr[0].user]) //当用户人数少于等于4，找相似度最高的1人
			} else {
				//当用户人数多于4人，找相似度最高的2人
				CurrentVoteMap[CurrentPostId] = CurrentAve + (sim[CurrentUserID][simStructArr[0].user]*(recommendMap[simStructArr[0].user][CurrentPostId]-averageArr[simStructArr[0].user])+sim[CurrentUserID][simStructArr[1].user]*(recommendMap[simStructArr[1].user][CurrentPostId]-averageArr[simStructArr[1].user]))/(sim[CurrentUserID][simStructArr[0].user]+sim[CurrentUserID][simStructArr[1].user])
			}
			//fmt.Printf("last %v %v %f\n", simStructArr[0].user, CurrentPostId, CurrentVoteMap[CurrentPostId])
			reCommendIdList = append(reCommendIdList, reCommendIdListStruct{CurrentPostId, CurrentVoteMap[CurrentPostId]})
		}
	}
	sort.Slice(reCommendIdList, func(i, j int) bool {
		return reCommendIdList[i].simVote > reCommendIdList[j].simVote
	})
	for _, listStruct := range reCommendIdList {
		split := strings.Split(listStruct.postKey, ":")
		reCommendKeyList = append(reCommendKeyList, split[3])
	}
	reCommendList := redis.GetReCommendList(reCommendKeyList)
	ResponseSuccess(c, reCommendList)
}
