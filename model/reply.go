package model

type Reply struct {
	Universal
	ReplyTo   int  `json:"-" gorm:"column:reply_to"`                         // 回复user_id
	ReplyUser User `gorm:"foreignKey:Id;references:ReplyTo" json:"reply_to"` //关联User
	CommentID int  `json:"comment_id" `
	User      User `gorm:"foreignKey:Id;references:UserID" json:"user"` //关联User
}

func (m *Reply) TableName() string {
	return "reply"
}

//type Reply struct {
//	CommentID int    `json:"comment_id" gorm:"column:comment_id"` // 评论id
//	ID        int    `json:"id" gorm:"column:id"`
//	UserID    string `json:"user_id" gorm:"column:user_id"`
//	Like      int    `json:"like" gorm:"column:like"`
//	Unlike    int    `json:"unlike" gorm:"column:unlike"`
//	Content   string `json:"content" gorm:"column:content"`
//	ReplyTo   string `json:"reply_to" gorm:"column:reply_to"` // 回复user_id
//}
//
//func (m *Reply) TableName() string {
//	return "reply"
//}
