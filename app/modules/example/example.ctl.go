package example

import (
	"app/app/modules/base"
	"app/app/modules/example/dto"
	"time"

	"github.com/gin-gonic/gin"
)

type ExampleCTL struct {
	service *ExampleService
}

func newCTL(service *ExampleService) *ExampleCTL {
	return &ExampleCTL{
		service: service,
	}
}

// https://gin-gonic.com/docs/examples/custom-validators/
func (ctl *ExampleCTL) User(ctx *gin.Context) {
	base.Success(ctx, ctl.service.User())
}

func (ctl *ExampleCTL) JSON(ctx *gin.Context) {
	base.JSON(ctx, 200, dto.ExampleDTOResponse{
		Name:      "test",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
}
