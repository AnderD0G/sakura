package model

import "sakura/pkg"

type Universal struct {
	ID          int    `json:"id" gorm:"column:id"`                     // 非业务主键
	Like        int    `json:"like" gorm:"column:like"`                 // 点赞数-公共属性
	Unlike      int    `json:"unlike" gorm:"column:unlike"`             // 不喜欢-公共属性
	UserID      string `json:"user_id" gorm:"column:user_id"`           // user_open_id-公共属性
	PublishTime string `json:"publish_time" gorm:"column:publish_time"` // 发布时间-公共属性
	Content     string `json:"content" gorm:"column:content"`           // 评论内容-公共内容

}

type Comment struct {
	Universal
	RelationType int     `json:"relation_type"` // 有哪些类型，暂时只有剧本
	RelationID   int     `json:"relation_id"`   // 联系的实体的id
	Replys       []Reply `gorm:"foreignKey:CommentID"`
}

type Reply struct {
	Universal
	ReplyTo   string `json:"reply_to" ` // 回复user_id
	CommentID int    `json:"comment_id" `
}

func (m *Reply) TableName() string {
	return "reply"
}

func (m *Comment) TableName() string {
	return "comment"
}

func GetComments(s *pkg.Inquirer[Comment]) {

}
