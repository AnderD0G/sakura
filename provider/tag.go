package provider

import (
	"context"
	"sakura/db"
	"sakura/model"
	"sakura/pkg"
)

type Tag struct {
	queryMap *pkg.QueryCondition
	tagModel *[]model.Tag
}

func (w *Tag) Work(ctx context.Context, finishChan chan<- pkg.Finish) {
	go pkg.Watcher(ctx, finishChan)

	d := db.GetMysql("1")

	t := new([]model.Tag)
	if err := d.Where("type = ?", "script").Find(t).Error; err != nil {
		pkg.SafeSend(finishChan, pkg.Finish{
			IsDone: false,
			Err:    err,
		})
	}

	w.tagModel = t
	pkg.SafeSend(finishChan, pkg.Finish{
		IsDone: true,
		Err:    nil,
	})

}
