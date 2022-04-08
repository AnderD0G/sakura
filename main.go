package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"sakura/db"
	"sakura/model"
	"sakura/pkg"
	"sakura/provider"
)

func init() {
	d := db.DB[*gorm.DB]{}

	d.Provider = &db.MysqlPro{Address: "root:123456qwe@tcp(127.0.0.1:3306)/taihe"}

	d.Initial()

	db.SetMysql(&d)
}
func main() {

	t := new(provider.HTTPHandler[model.Scripts])
	t.Provider = &provider.Scripts{QueryMap: new(pkg.QueryCondition)}

	router := gin.Default()

	router.GET("/scripts", t.List())
	log.Fatal(router.Run(":8080"))

}
