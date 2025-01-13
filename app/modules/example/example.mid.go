package example

type ExampleMiddleware struct {
	svc *ExampleService
}

func newMiddleware(svc *ExampleService) *ExampleMiddleware {
	return &ExampleMiddleware{
		svc: svc,
	}
}
