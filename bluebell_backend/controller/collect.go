package controller

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type CollectStruct struct {
	PostID    string  `json:"post_id"`
	Direction float64 `json:"direction"`
}

func CollectHandler(c *gin.Context) {
	var collect CollectStruct
	if err := c.ShouldBindJSON(&collect); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	userId, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	if collect.Direction == 1 {
		if err = mysql.Collect(collect.PostID, userId); err != nil {
			zap.L().Error("controller.CollectHandler() failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	} else if collect.Direction == 0 {
		if err = mysql.CancelCollect(collect.PostID, userId); err != nil {
			zap.L().Error("controller.CollectHandler() failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	}
	ResponseSuccess(c, nil)
}

func CollectListHandler(c *gin.Context) {
	pageStr, _ := c.GetQuery("page")
	pageNum, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	userId, _ := GetCurrentUserID(c)
	collectList := mysql.GetCollectList(userId)
	fmt.Printf("%v\n", collectList)
	postList := redis.GetPostBYKeys(collectList, pageNum)
	ResponseSuccess(c, postList)
}
