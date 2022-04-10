package provider

import "C"
import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"sakura/model"
	"sakura/pkg"
	"time"
)

type Scripts struct {
	QueryMap    *pkg.QueryCondition
	scriptModel *[]model.Script
	S           *pkg.Inquirer[*model.Script]
}

func (t *Scripts) FindByID(context *gin.Context) (model.Script, error) {
	panic("implement me")
	//return model.Script{Name: "luiz"}, nil
}

func (s *Scripts) List(c *gin.Context) ([]model.Script, error) {

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	query := c.DefaultQuery("query", "")

	s.QueryMap.Query = query
	s.QueryMap.Page = pkg.Atoi(page)
	s.QueryMap.Size = pkg.Atoi(size)

	t := new(Tag)

	err := pkg.Run(2*time.Second, c, s, t)
	if err != nil {
		return nil, err
	}

	m := make(map[int]string)
	for _, v := range *t.tagModel {
		m[v.Uuid] = v.Value
	}

	scripts := make([]model.Script, len(*s.scriptModel))

	for k, v := range *s.scriptModel {
		tags := make([]string, len(v.ScriptTag))
		for k, v := range v.ScriptTag {
			if s, ok := m[v]; ok {
				tags[k] = s
			} else {
				return nil, errors.New(fmt.Sprintf("tag %v 没有对应value", k))
			}
		}
		v.Tags = tags
		scripts[k] = v
	}

	return scripts, err
}

func (t *Scripts) Update(id string, model model.Script) error {
	//TODO implement me
	panic("implement me")
}

func (t *Scripts) Insert(model model.Script) error {
	//TODO implement me
	panic("implement me")
}

func (t *Scripts) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Scripts) Work(ctx context.Context, finishChan chan<- pkg.Finish) {
	go pkg.Watcher(ctx, finishChan)
	i := new([]model.Script)

	s.S.GetParam(s.QueryMap)
	s.S.ParseStruct()
	err := s.S.ParseRule()
	if err != nil {
		pkg.SafeSend(finishChan, pkg.Finish{
			IsDone: false,
			Err:    err,
		})
	}

	s.S.Query(i)

	s.scriptModel = i
	pkg.SafeSend(finishChan, pkg.Finish{
		IsDone: true,
		Err:    nil,
	})
}
