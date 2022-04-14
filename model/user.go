package model

import (
	"context"
	db2 "sakura/db"
)

type (
	UserPub struct {
		Id        int    `json:"-" gorm:"column:id"`
		NickName  string `json:"nick_name" gorm:"column:nick_name"`   // 用户昵称
		AvatarUrl string `json:"avatar_url" gorm:"column:avatar_url"` // 用户头像图片的 URL。URL 最后一个数值代表正方形头像大小（有 0、46、64、96、132 数值可选，0 代表 640x640 的正方形头像，46 表示 46x46 的正方形头像，剩余数值以此类推。默认132），用户没有头像时该项为空。若用户更换头像，原有头像 URL 将失效。
		Gender    int    `json:"gender" gorm:"column:gender"`
	}
	User struct {
		UserPub
		Openid     string `json:"-" gorm:"column:openid"`      // 用户唯一标识
		SessionKey string `json:"-" gorm:"column:session_key"` // 会话密钥
		TypeId     string `json:"-" gorm:"column:typeid"`      //  对应的某个小程序

	}
)

func (m *User) TableName() string {
	return "user"
}

func GetUserInfo(ctx context.Context, typeId, openId string) (*User, error) {
	db := db2.GetMysql("1")
	i := new(User)
	affected := db.First(i, "typeid = ? AND openid = ?", typeId, openId).RowsAffected
	if affected == 1 {
		return i, nil
	}
	return nil, nil
}
