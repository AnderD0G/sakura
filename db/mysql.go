package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *DB[*gorm.DB]

type MysqlPro struct {
	m       map[string]*gorm.DB
	Address string
}

func (g *MysqlPro) initial() error {
	db, err := gorm.Open(mysql.Open(g.Address), &gorm.Config{})
	if err != nil {
		return err
	}
	g.m = make(map[string]*gorm.DB)
	g.m["1"] = db
	return nil

}

func (g *MysqlPro) instance(key string) *gorm.DB {
	return g.m[key]
}

func SetMysql(d *DB[*gorm.DB]) {
	db = d
}

func GetMysql(key string) *gorm.DB {
	return db.Instance(key)
}
