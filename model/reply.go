package model

type Reply struct {
	Universal
	ReplyTo   string `json:"reply_to" ` // 回复user_id
	CommentID int    `json:"comment_id" `
}

func (m *Reply) TableName() string {
	return "reply"
}
