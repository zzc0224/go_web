package models

import (
	"encoding/json"
	"errors"
	"time"
)

type Post struct {
	FileList    string    `json:"fileList" db:"images"`
	PostID      uint64    `json:"post_id,string" db:"post_id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	AuthorId    uint64    `json:"author_id,string" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"-" db:"create_time"`
}

func (p *Post) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Title       string `json:"title" db:"title"`
		Content     string `json:"content" db:"content"`
		CommunityID int64  `json:"community_id" db:"community_id"`
		FileList    string `json:"fileList" db:"images"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.Title) == 0 {
		err = errors.New("帖子标题不能为空")
	} else if len(required.Content) == 0 {
		err = errors.New("帖子内容不能为空")
	} else if required.CommunityID == 0 {
		err = errors.New("未指定版块")
	} else {
		p.Title = required.Title
		p.Content = required.Content
		p.CommunityID = required.CommunityID
		p.FileList = required.FileList
	}
	return
}

type ApiPostDetail struct {
	*Post
	AuthorName    string  `json:"author_name"`
	CommunityName string  `json:"community_name"`
	Direction     float64 `json:"direction"`
	IsConcern     float64 `json:"is_concern"`
}
