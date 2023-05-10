package controller

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	//page, _ := c.GetQuery("page")
	//
	//type structOrder struct {
	//	Order string `json:"keywords"`
	//}
	//var O structOrder
	//if err := c.ShouldBindJSON(&O); err != nil {
	//	zap.L().Error("invalid params", zap.Error(err))
	//	ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
	//	return
	//}
	//println(O.Order)
	keywords := c.Param("keywords")
	page, _ := c.GetQuery("page")
	postIDList := mysql.GetPostBYOrder(keywords)
	keys := make([]string, 0)
	for _, s := range postIDList {
		keys = append(keys, strconv.FormatUint(s, 10))
	}
	pageNum, _ := strconv.ParseInt(page, 10, 64)
	postList := redis.GetPostBYKeys(keys, pageNum)
	ResponseSuccess(c, postList)
}
