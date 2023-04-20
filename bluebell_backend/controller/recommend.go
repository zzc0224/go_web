package controller

import (
	"bluebell_backend/dao/redis"

	"github.com/gin-gonic/gin"
)

func ReCommend(c *gin.Context) {
	recommendMap := redis.GetVote()
	for userId, List := range recommendMap {
		for postId, vote := range List {
			println(userId, postId, vote)
		}
	}
}
