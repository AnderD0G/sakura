package provider

import (
	"github.com/gin-gonic/gin"
	"go/types"
	"gorm.io/gorm"
	"net/http"
	"sakura/model"
)

type HTTPHandler[MODEL model.Model] struct {
	db              *gorm.DB
	Provider        Provider[MODEL]
	InsertValidator func(new MODEL) error
}

func (h *HTTPHandler[MODEL]) BindDB(db *gorm.DB) {
	h.db = db
}

func (h *HTTPHandler[MODEL]) FindByID() gin.HandlerFunc {

	return func(context *gin.Context) {
		if r, err := h.Provider.FindByID(context); err == nil {
			context.String(http.StatusOK, "%v", r)
		}
	}
}

func (h *HTTPHandler[MODEL]) List() gin.HandlerFunc {

	return func(context *gin.Context) {
		if r, err := h.Provider.List(context); err == nil {
			context.JSON(http.StatusOK, r)
		}
	}
}

type Samples interface {
	*gorm.DB | types.Nil
}

type HTTPHandlers[MODEL Samples] struct {
}
