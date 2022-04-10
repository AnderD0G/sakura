package provider

import (
	"github.com/gin-gonic/gin"
	"sakura/model"
	"sakura/pkg"
)

type Comment struct {
	QueryMap *pkg.QueryCondition
	I        *pkg.Inquirer[model.Comment]
}

func (c *Comment) FindByID(context *gin.Context) (model.Comment, error) {

	panic("implement me")
	//return model.Script{Name: "luiz"}, nil
}

func (s *Comment) List(c *gin.Context) ([]model.Comment, error) {

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	query := c.DefaultQuery("query", "")

	s.QueryMap.Query = query
	s.QueryMap.Page = pkg.Atoi(page)
	s.QueryMap.Size = pkg.Atoi(size)

	s.I.GetParam(s.QueryMap)
	s.I.ParseStruct()
	if err := s.I.ParseRule(); err != nil {
		return nil, err
	}
	//s.I.Query()
	return nil, nil
}

func (t *Comment) Update(id string, model model.Comment) error {
	//TODO implement me
	panic("implement me")
}

func (t *Comment) Insert(model model.Comment) error {
	//TODO implement me
	panic("implement me")
}

func (t *Comment) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
