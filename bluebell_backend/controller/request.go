package controller

import (
	"bluebell_backend/dao/mysql"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserID(c *gin.Context) (userID uint64, err error) {
	_userID, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = _userID.(uint64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func GetCurrUserName(c *gin.Context) {
	currentUserID, _ := GetCurrentUserID(c)
	user, _ := mysql.GetUserByID(strconv.FormatUint(currentUserID, 10))
	ResponseSuccess(c, user)
}
