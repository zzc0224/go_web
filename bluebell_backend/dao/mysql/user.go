package mysql

import (
	"bluebell_backend/models"
	"bluebell_backend/pkg/snowflake"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"strings"
)

const secret = "liwenzhou.com"

func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Register(user *models.User) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int64
	err = db.Get(&count, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count > 0 {
		// 用户已存在
		return ErrorUserExit
	}
	// 生成user_id
	userID, err := snowflake.GetID()
	if err != nil {
		return ErrorGenIDFailed
	}
	// 生成加密密码
	password := encryptPassword([]byte(user.Password))
	// 把用户插入数据库
	sqlStr = "insert into user(user_id, username, password) values (?,?,?)"
	_, err = db.Exec(sqlStr, userID, user.UserName, password)
	return
}

func Login(user *models.User) (err error) {
	originPassword := user.Password // 记录一下原始密码
	sqlStr := "select user_id, username, password,status from user where username = ?"
	err = db.Get(user, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		return
	}
	if err == sql.ErrNoRows {
		// 用户不存在
		return ErrorUserNotExit
	}
	// 生成加密密码与查询到的密码比较
	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return ErrorPasswordWrong
	}
	if user.Status != 0 {
		return ErrorUserProhibit
	}
	return
}

func GetUserByID(idStr string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, idStr)
	return
}

func GetAllUser() []string {
	var userId []string
	sqlStr := `select user_id from user`
	db.Select(&userId, sqlStr)
	return userId
}

func GetCollectList(userId uint64) []string {
	var collectList string
	sqlStr := `select collect from user where user_id = ?`
	db.Get(&collectList, sqlStr, userId)
	split := strings.Split(collectList, ";")
	split = append(split[:0], split[1:]...)
	return split
}

func GetConcernList(userId uint64) []string {
	var concernList string
	sqlStr := `select concern from user where user_id = ?`
	db.Get(&concernList, sqlStr, userId)
	split := strings.Split(concernList, ";")
	split = append(split[:0], split[1:]...)
	return split
}

func Collect(postId string, userId uint64) (err error) {
	var collectList string
	sqlStr := `select collect from user where user_id = ?`
	db.Get(&collectList, sqlStr, userId)
	collectList = collectList + ";" + postId
	sqlStr = `update user set collect = ? where user_id = ?`
	_, err = db.Exec(sqlStr, collectList, userId)
	return err
}

func CancelCollect(postId string, userId uint64) (err error) {
	collectList := GetCollectList(userId)
	var collectPostIdList string
	for i := 0; i < len(collectList); i++ {
		if collectList[i] != postId && collectList[i] != "" {
			collectPostIdList = collectPostIdList + ";" + collectList[i]
		}
	}
	sqlStr := `update user set collect = ? where user_id = ?`
	_, err = db.Exec(sqlStr, collectPostIdList, userId)
	return err
}

func Concern(concernUserid string, userId uint64) (err error) {
	var concernList string
	sqlStr := `select concern from user where user_id = ?`
	db.Get(&concernList, sqlStr, userId)
	concernList = concernList + ";" + concernUserid
	sqlStr = `update user set concern = ? where user_id = ?`
	_, err = db.Exec(sqlStr, concernList, userId)
	return
}

func CancelConcern(concernUserid string, userId uint64) (err error) {
	concernList := GetConcernList(userId)
	var concernUserIdList string
	for i := 0; i < len(concernList); i++ {
		if concernList[i] != concernUserid && concernList[i] != "" {
			concernUserIdList = concernUserIdList + ";" + concernList[i]
		}
	}
	sqlStr := `update user set concern = ? where user_id =?`
	_, err = db.Exec(sqlStr, concernUserIdList, userId)
	return
}
