package model

type Shop struct {
	ShopName   string  `json:"shop_name" gorm:"column:shop_name"`
	ShopAddr   string  `json:"shop_addr" gorm:"column:shop_addr"`
	ShopScore  float64 `json:"shop_score" gorm:"column:shop_score"`
	Numone     string  `json:"numone" gorm:"column:numone"`
	Numtwo     string  `json:"numtwo" gorm:"column:numtwo"`
	ShopLogo   string  `json:"shop_logo" gorm:"column:shop_logo"`
	ShopCover  string  `json:"shop_cover" gorm:"column:shop_cover"`
	MID        string  `json:"m_id" gorm:"column:m_id"`
	Longtitude float64 `json:"longtitude" gorm:"column:longtitude"`
	Latitude   float64 `json:"latitude" gorm:"column:latitude"`
	ID         int     `json:"id" gorm:"column:id"`
}

func (m *Shop) TableName() string {
	return "shop"
}
