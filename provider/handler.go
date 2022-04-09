package provider

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sakura/model"
)

type (
	Provider[MODEL model.Model] interface {
		FindByID(context *gin.Context) (MODEL, error)
		List(context *gin.Context) ([]MODEL, error)
		Update(id string, model MODEL) error
		Insert(model MODEL) error
		Delete(id string) error
	}

	HTTPHandler[MODEL model.Model] struct {
		Provider        Provider[MODEL]
		InsertValidator func(new MODEL) error
	}
)

func (h *HTTPHandler[MODEL]) List() gin.HandlerFunc {

	return func(context *gin.Context) {
		if r, err := h.Provider.List(context); err == nil {
			context.JSON(http.StatusOK, r)
		}
	}
}

func (h *HTTPHandler[MODEL]) FindByID() gin.HandlerFunc {

	return func(context *gin.Context) {
		if r, err := h.Provider.FindByID(context); err == nil {
			context.JSON(http.StatusOK, r)
		}
	}
}
