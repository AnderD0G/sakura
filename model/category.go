package model

type (
	Category struct {
		Value string `json:"value" gorm:"column:value"`
		Uuid  string `json:"uuid" gorm:"primary_key" valid:"no_empty"`
		IsDel int    `json:"-" gorm:"column:is_del"`
		Tags  *[]Tag `json:"tags,omitempty" gorm:"many2many:c_tag;foreignKey:Uuid;joinForeignKey:CategoryUid;References:Uuid;JoinReferences:TagUid" valid:"no_empty"`
	}
	CategoriesTagDto struct {
		TIds []string `json:"tag_ids" valid:"no_empty"`
		Cid  string   `json:"category_id" valid:"no_empty"`
	}
)
