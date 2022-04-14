package model

import (
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"sakura/pkg"
)

type (
	UniversalPub struct {
		Like        int    `json:"like" gorm:"column:like"`                 // 点赞数-公共属性
		PublishTime string `json:"publish_time" gorm:"column:publish_time"` // 发布时间-公共属性
		Content     string `json:"content" gorm:"column:content"`           // 评论内容-公共内容
		ID          int    `json:"id" gorm:"column:id"`                     // 非业务主键
	}
	Universal struct {
		UniversalPub
		Unlike int    `json:"unlike" gorm:"column:unlike"`   // 不喜欢-公共属性
		UserID string `json:"user_id" gorm:"column:user_id"` // user_open_id-公共属性
	}
)
type (
	CommentDetail struct {
		Comment
		Replies []Reply `gorm:"foreignKey:CommentID" json:"replies"`
	}
	Comment struct {
		Universal
		RelationType int                                                   `json:"relation_type"` // 有哪些类型，暂时只有剧本
		RelationID   string                                                `json:"relation_id"`   // 联系的实体的id
		User         `gorm:"foreignKey:Id;references:UserID;" json:"user"` //关联User
		ReplyCounts  int                                                   `json:"reply_counts" gorm:"column:reply_counts"`
	}
	CommentPub struct {
		UniversalPub
		UserPub     `json:"user"`
		ReplyCounts int `json:"reply_counts" gorm:"column:reply_counts"`
	}
)

func (m *Comment) TableName() string {
	return "comment"
}

func GetComments(s *pkg.Inquirer[*Comment]) []Comment {
	i := make([]Comment, 0)
	k := func(db *gorm.DB) {
		s.Db.Table("(?)as u", db).Select("*").Joins("left join user on u.user_id = user.id").Order("publish_time desc").Find(&i)
	}
	s.Query("comment", nil, k)
	return i
}

func GetComment(s *pkg.Inquirer[*Comment]) Comment {
	return GetComments(s)[0]
}

func CommentsPub(from *[]Comment) (error, interface{}) {
	pubs := make([]CommentPub, 0)
	err := copier.Copy(&pubs, from)
	if err != nil {
		return errors.New(http.StatusText(http.StatusInternalServerError)), nil
	}
	return nil, pubs
}
