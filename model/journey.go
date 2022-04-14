package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	db2 "sakura/db"
)

type (
	PersonBasicInfo struct {
		Id     string `json:"id"`
		Gender int    `json:"gender"`
	}
	Persons []PersonBasicInfo
)

func (p Persons) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Persons) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &p)
}

type JourneyDis struct {
	Journey
	Persons Persons `json:"personp" gorm:"column:personp"`
	Dis     float64 `json:"dis"`
}

type Journey struct {
	Time     string  `json:"time" gorm:"column:time"`
	Title    string  `json:"title" gorm:"column:title"`
	ID       int     `json:"id" gorm:"column:id"`
	Shop     Shop    `json:"shop" gorm:"foreignKey:MID;references:ShopID"`
	Scripts  Script  `json:"script" gorm:"foreignKey:Uuid;references:ScriptID"`
	Persons  Persons `json:"-" gorm:"column:persons"`
	ScriptID string  `json:"script_id" gorm:"column:script_id"`
	Owner    string  `json:"owner" gorm:"column:owner"`
	ShopID   string  `json:"shop_id" gorm:"column:shop_id"`
	Status   int     `json:"status" gorm:"column:status"` // 1:进行中 2:已开始 3:已结束
}

func (m *Journey) TableName() string {
	return "journey"
}

func GetJourney(long, lat string) []JourneyDis {
	j := make([]JourneyDis, 0)

	db := db2.GetMysql("1")
	sub := db.Model(&Journey{}).Joins("Shop").Joins("Scripts").Where("status = ?", 1)
	db.Table("(?)as u", sub).Select(fmt.Sprintf("*,round(st_distance_sphere(point(%v,%v),point(Shop__longtitude,Shop__latitude))) dis,@age:= CONCAT('$[0 to ',Scripts__script_player_limit-1,' ]'),JSON_EXTRACT(persons, @age)as personp", long, lat)).Order("dis asc").Find(&j)

	return j
}

type JourneyPerson struct {
	Journey
	Persons []User `json:"personp" gorm:"column:personp"`
	Dis     int    `json:"dis"`
}

func GetJourneyDetailM(id int) (*JourneyPerson, error) {
	db := db2.GetMysql("1")
	j := new(JourneyDis)

	sub := db.Model(&Journey{ID: id}).Joins("Shop").Joins("Script")
	db.Table("(?)as u", sub).Select(fmt.Sprintf("*,@age:= CONCAT('$[0 to ',Scripts__script_player_limit-1,' ]'),JSON_EXTRACT(persons, @age)as personp")).Find(&j)

	i := new([]User)
	strings := make([]string, len(j.Persons))

	err := copier.Copy(&strings, j.Persons)
	if err != nil {
		return nil, err
	}

	if err = db.Where("uuid IN ?", strings).Find(i).Error; err != nil {
		return nil, err
	}

	j2 := new(JourneyPerson)
	j2.Journey = j.Journey
	j2.Persons = *i

	return j2, nil
}
