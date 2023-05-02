package controller

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	type structOrder struct {
		Order string `json:"order"`
	}
	var O structOrder
	if err := c.ShouldBindJSON(&O); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	println(O.Order)
	postIDList := mysql.GetPostBYOrder(O.Order)
	keys := make([]string, 0)
	for _, s := range postIDList {
		keys = append(keys, strconv.FormatUint(s, 10))
	}
	postList := redis.GetPostBYKeys(keys)
	ResponseSuccess(c, postList)
}
