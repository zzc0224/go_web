package controller

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	UserPerAge = 3
)

type ConcernStruct struct {
	ConcernUserID string  `json:"author_id"`
	IsConcern     float64 `json:"is_concern"`
}

type ConcernUserStruct struct {
	ConcernUser []*models.User `json:"concern_user"`
	ConcernSum  int64          `json:"concern_sum""`
}

func ConcernHandler(c *gin.Context) {
	var concern ConcernStruct
	if err := c.ShouldBindJSON(&concern); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	userId, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	if concern.IsConcern == 1 {
		if err = mysql.Concern(concern.ConcernUserID, userId); err != nil {
			zap.L().Error("controller.ConcernHandler() failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	} else if concern.IsConcern == 0 {
		if err = mysql.CancelConcern(concern.ConcernUserID, userId); err != nil {
			zap.L().Error("controller.ConcernHandler() failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	}
	ResponseSuccess(c, nil)
}

func ConcernListHandler(c *gin.Context) {
	var concernUser ConcernUserStruct
	pageStr, _ := c.GetQuery("page")
	pageNum, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	start := (pageNum - 1) * UserPerAge
	end := start + UserPerAge - 1
	userId, _ := GetCurrentUserID(c)
	concernList := mysql.GetConcernList(userId)
	concernUser.ConcernSum = int64(len(concernList))
	var concernUserList []*models.User
	for i := start; i <= end; i++ {
		if i < concernUser.ConcernSum {
			user, _ := mysql.GetUserByID(concernList[i])
			concernUserList = append(concernUserList, user)
		} else {
			break
		}
	}
	concernUser.ConcernUser = concernUserList

	ResponseSuccess(c, concernUser)
}
