package provider

import (
	"context"
	"github.com/gin-gonic/gin"
	"sakura/model"
	"sakura/pkg"
	"time"
)

type Comment struct {
	Query *pkg.Query
	I     *pkg.Inquirer[*model.Comment]
	//仅用作保存worker返回结果
	c *model.Comment
	//存放Reply查询器，用作查出关联的数据
	R *pkg.Inquirer[*model.Reply]
}

func (c *Comment) FindByID(context *gin.Context) (model.Comment, error) {

	query := context.DefaultQuery("query", "")
	c.Query.Condition = query
	c.Query.Condition = ""
	err := pkg.Run(2*time.Second, context, c)
	if err != nil {
		return model.Comment{}, err
	}
	return *c.c, nil
}

func (s *Comment) List(c *gin.Context) ([]model.Comment, error) {

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	query := c.DefaultQuery("query", "")

	s.Query.Condition = query
	s.Query.Page = pkg.Ati(page)
	s.Query.Size = pkg.Ati(size)

	s.I.InjectParam(s.Query)
	s.I.ParseStruct()

	if err := s.I.ParseQuery(); err != nil {
		return nil, err
	}

	comments := model.GetComments(s.I)
	return comments, nil
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

func (w *Comment) Work(ctx context.Context, finishChan chan<- pkg.Finish) {
	go pkg.Watcher(ctx, finishChan)

	w.I.InjectParam(w.Query)
	w.I.ParseStruct()

	if err := w.I.ParseQuery(); err != nil {
		pkg.SafeSend(finishChan, pkg.Finish{
			IsDone: false,
			Err:    err,
		})
	}
	if err := w.I.ParseQuery(); err != nil {
		pkg.SafeSend(finishChan, pkg.Finish{
			IsDone: false,
			Err:    err,
		})
	}

	comment := model.GetComment(w.I)
	w.c = &comment
	pkg.SafeSend(finishChan, pkg.Finish{
		IsDone: true,
		Err:    nil,
	})

}
