package mysql

import (
	"bluebell_backend/models"

	"go.uber.org/zap"
)

func CreateComment(comment *models.Comment) (err error) {
	sqlStr := `insert into comment(
	comment_id, content, post_id, author_id, parent_id,author_name)
	values(?,?,?,?,?,?)`
	_, err = db.Exec(sqlStr, comment.CommentID, comment.Content, comment.PostID,
		comment.AuthorID, comment.ParentID, comment.AuthorName)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func GetCommentListByIDs(ids string) (commentList []*models.Comment, err error) {
	sqlStr := `select comment_id, content, post_id, author_id, parent_id, create_time,author_name
	from comment
	where post_id = ?`
	commentList = make([]*models.Comment, 0, 2)
	err = db.Select(&commentList, sqlStr, ids)
	//// 动态填充id
	//query, args, err := sqlx.In(sqlStr, ids)
	//if err != nil {
	//	return
	//}
	//// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	//query = db.Rebind(query)
	//err = db.Select(&commentList, query, args...)
	return
}

func GetCommentAuthorName(userID uint64) (AuthorName string, err error) {
	sqlStr := `select username
	from user
	where user_id = ?`
	row := db.QueryRow(sqlStr, userID)
	row.Scan(&AuthorName)
	return
}
