package mysql

import (
	"bluebell_backend/models"

	"go.uber.org/zap"
)

func CreateWarn(warn *models.Warn) (err error) {
	sqlStr := `insert into warn(warning_id,positive_user_id,passive_user_id,comment_post_user_id,type,status,reason)
	values(?,?,?,?,?,?,?)`
	_, err = db.Exec(sqlStr, warn.WarnID, warn.PositiveUserID, warn.PassiveUserID, warn.CommentPostUserID, warn.Type, warn.Status, warn.Reason)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func GetWarnList() (WarnList []*models.Warn) {
	sqlStr := `select warning_id,positive_user_id,passive_user_id,comment_post_user_id,type,reason from warn where status = 0`
	//sqlStr := `select * from warn`
	WarnList = make([]*models.Warn, 0)
	db.Select(&WarnList, sqlStr)
	return
}

func AdminWarn(warn *models.Warn) {
	sqlStr := `update warn set status = ?,positive_result = ?,passive_result = ? where warning_id = ?`
	_, err := db.Exec(sqlStr, warn.Status, warn.PositiveResult, warn.PassiveResult, warn.WarnID)
	if err != nil {
		zap.L().Error("update warn failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
}

func DeleteComment(commentId uint64) {
	sqlStr := `delete from comment where comment_id = ?`
	_, err := db.Exec(sqlStr, commentId)
	if err != nil {
		zap.L().Error("delete comment failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
}

func DeletePost(PostId uint64) {
	sqlStr := `delete from post where post_id = ?`
	_, err := db.Exec(sqlStr, PostId)
	if err != nil {
		zap.L().Error("delete post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
}

func AdminUser(UserId uint64) {
	sqlStr := `update user set status = ? where user_id = ?`
	_, err := db.Exec(sqlStr, 1, UserId)
	if err != nil {
		zap.L().Error("delete post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
}
