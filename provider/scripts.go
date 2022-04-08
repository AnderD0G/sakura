package provider

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	db2 "sakura/db"
	"sakura/model"
	"sakura/pkg"
	"time"
)

type Scripts struct {
	QueryMap    *pkg.QueryCondition
	scriptModel *[]model.Scripts
	total       int64
}

type Provider[MODEL model.Model] interface {
	FindByID(context *gin.Context) (MODEL, error)
	List(context *gin.Context) ([]MODEL, error)
	Update(id string, model MODEL) error
	Insert(model MODEL) error
	Delete(id string) error
}

func (t *Scripts) FindByID(context *gin.Context) (model.Scripts, error) {
	panic("implement me")
	//return model.Scripts{Name: "luiz"}, nil
}

func (s *Scripts) List(c *gin.Context) ([]model.Scripts, error) {

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	query := c.DefaultQuery("query", "")

	s.QueryMap.Query = query
	s.QueryMap.Page = int64(pkg.Atoi(page))
	s.QueryMap.Size = int64(pkg.Atoi(size))

	t := new(Tag)
	err := pkg.Run(2*time.Second, c, s, t)
	if err != nil {
		return nil, err
	}

	m := make(map[int]string)
	for _, v := range *t.tagModel {
		m[v.Uuid] = v.Value
	}

	scripts := make([]model.Scripts, len(*s.scriptModel))

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

func (t *Scripts) Update(id string, model model.Scripts) error {
	//TODO implement me
	panic("implement me")
}

func (t *Scripts) Insert(model model.Scripts) error {
	//TODO implement me
	panic("implement me")
}

func (t *Scripts) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Scripts) Work(ctx context.Context, finishChan chan<- pkg.Finish) {
	go pkg.Watcher(ctx, finishChan)
	num := int64(0)
	i := new([]model.Scripts)

	db := db2.GetMysql("1")

	page, limit, query := pkg.GetParam(s.QueryMap)
	queryMap, err := pkg.GenerateKv(query)

	if err != nil {
		pkg.SafeSend(finishChan, pkg.Finish{
			IsDone: false,
			Err:    err,
		})
	}

	for k, v := range queryMap {
		if v.Rule == pkg.Normal {
			db = db.Where(fmt.Sprintf("%v %v ?", k, v.Comparator), v.Value[0])
		}
		if v.Rule == pkg.JsonArray {
			db = db.Where(v.Value[0])
		}
		if v.Rule == pkg.Array {
			db = db.Where(fmt.Sprintf("%v IN ?", k), v.Value)
		}
	}

	if num = db.Find(i).RowsAffected; num < 0 {
		pkg.SafeSend(finishChan, pkg.Finish{
			IsDone: false,
			Err:    errors.New("RowsAffected < 0"),
		})
	}

	if limit > 0 {
		db = db.Limit(limit).Offset((page - 1) * limit)
	}

	db.Find(i)

	s.scriptModel = i
	s.total = num
	pkg.SafeSend(finishChan, pkg.Finish{
		IsDone: true,
		Err:    nil,
	})
}
