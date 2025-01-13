package example

import (
	"app/app/modules/user"
)

type ExampleModule struct {
	Ctl *ExampleCTL
	Svc *ExampleService
	Mid *ExampleMiddleware
}

func New(userService *user.UserService) *ExampleModule {
	service := newService(userService)
	ctl := newCTL(service)
	mid := newMiddleware(service)
	return &ExampleModule{
		Ctl: ctl,
		Svc: service,
		Mid: mid,
	}
}
