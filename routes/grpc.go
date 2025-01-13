package routes

import (
	"app/app/modules"

	"google.golang.org/grpc"
)

// RegisterGRPCServer gRPC register service
func RegisterGRPCServer(s *grpc.Server, mod *modules.Modules) {
	// user.RegisterUserServer(s, controller.UserController)
}
