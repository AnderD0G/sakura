package model

type (
	Tag struct {
		Value string      `json:"value" gorm:"column:value"`
		Uuid  int         `json:"uuid" gorm:"primary_key" valid:"no_empty"`
		IsDel int         `json:"-" gorm:"column:is_del" value:"1|0"`
		Cs    *[]Category `json:"categories,omitempty" gorm:"many2many:c_tag;foreignKey:Uuid;joinForeignKey:TagUid;References:Uuid;JoinReferences:CategoryUid" valid:"no_empty"`
	}
)

func (m *Tag) TableName() string {
	return "tag"
}

type Model interface {
	Tag | Script | Category | JourneyDis | JourneyPerson | Comment
}
