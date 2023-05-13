package controller

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"bluebell_backend/models"
	"bluebell_backend/pkg/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WarnHandler(c *gin.Context) {
	var w models.Warn
	if err := c.ShouldBindJSON(&w); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	userID, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	w.PositiveUserID = userID
	warningID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}
	w.WarnID = warningID
	w.Status = "0"
	if err := mysql.CreateWarn(&w); err != nil {
		zap.L().Error("mysql.CreateWarn(&w) failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	ResponseSuccess(c, nil)
}

func WarnListHandler(c *gin.Context) {
	warnList := mysql.GetWarnList()
	ResponseSuccess(c, warnList)
}

func AdminWarn(c *gin.Context) {
	var w models.Warn
	if err := c.ShouldBindJSON(&w); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	mysql.AdminWarn(&w)
	if w.Status == "2" { //status为2，举报有效需要删除;status为1，举报无效不需删除
		if w.Type == "1" {
			mysql.DeleteComment(w.CommentPostUserID)
		} else if w.Type == "2" {
			communityName := mysql.GetCommunityNameByPostID(w.CommentPostUserID)
			redis.DeletePost(w.CommentPostUserID, communityName)
			mysql.DeletePost(w.CommentPostUserID)
		}
	}
	ResponseSuccess(c, nil)
}

func AdminUser(c *gin.Context) {
	var w models.Warn
	if err := c.ShouldBindJSON(&w); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	mysql.AdminUser(w.PassiveUserID)
	ResponseSuccess(c, nil)
}
