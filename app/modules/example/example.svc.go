package example

import (
	"app/app/modules/example/dto"
	"app/app/modules/user"
	"time"

	"github.com/jinzhu/copier"
)

type ExampleService struct {
	userService *user.UserService
}

func newService(userService *user.UserService) *ExampleService {
	return &ExampleService{
		userService: userService,
	}
}

func (service *ExampleService) User() *dto.UserDTOResponse {
	user := service.userService.User()
	ret := new(dto.UserDTOResponse)
	copier.CopyWithOption(ret, user, copier.Option{IgnoreEmpty: true, Converters: []copier.TypeConverter{
		{
			SrcType: time.Time{},
			DstType: int64(0),
			Fn: func(src interface{}) (interface{}, error) {
				return src.(time.Time).Unix(), nil
			},
		},
	}})
	return ret
}
