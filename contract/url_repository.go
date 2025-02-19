package contract

import (
	"go_link_reducer/types"

	"github.com/gin-gonic/gin"
)

type URLRepository interface {
	Create(types.CreateURLPayload) (types.URL, error)
	Update(ID uint, count int) error
	GetAll(c *gin.Context) (map[string]any, error)
	GetOne(URL string) (types.URL, error)
	Delete() error
}
