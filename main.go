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
	tai := db.GetMysql("1")

	t := new(provider.HTTPHandler[model.Script])
	t.Provider = &provider.Scripts{QueryMap: new(pkg.QueryCondition), S: &pkg.Inquirer[*model.Script]{
		M:        new(model.Script),
		N:        make(pkg.TypeMap),
		QueryMap: make(map[string]pkg.RuleType),
		Db:       tai,
	}}

	//journeyList
	j := new(provider.HTTPHandler[model.JourneyDis])
	j.Provider = &provider.Journey{}

	//journeyDetail
	d := new(provider.HTTPHandler[model.JourneyPerson])
	d.Provider = &provider.Detail{}

	router := gin.Default()
	router.GET("/scripts", t.List())
	router.GET("/js", t.List())
	router.GET("/js/detail", d.FindByID())
	log.Fatal(router.Run(":8080"))

}
