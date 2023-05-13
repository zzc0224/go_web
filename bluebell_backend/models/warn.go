package models

import "time"

type Warn struct {
	WarnID            uint64    `json:"warn_id,string" db:"warning_id"`
	PositiveUserID    uint64    `json:"positive_user_id,string" db:"positive_user_id"`
	PassiveUserID     uint64    `json:"passive_user_id,string" db:"passive_user_id"`
	CommentPostUserID uint64    `json:"comment_post_user_id,string" db:"comment_post_user_id"`
	Type              string    `json:"type" db:"type"`
	Status            string    `json:"status" db:"status"`
	PositiveResult    string    `json:"positive_result" db:"positive_result"`
	PassiveResult     string    `json:"passive_result" db:"passive_result"`
	CreateTime        time.Time `json:"create_time" db:"create_time"`
	Reason            string    `json:"reason" db:"reason"`
}
