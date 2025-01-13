package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/app/modules"
	"app/internal/modules/log"
	"app/internal/ssl"
	"app/routes"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	// "google.golang.org/grpc/credentials"
)

// GRPC is gRPC server
func GRPC(isGRPCS bool) *cobra.Command {
	name := "grpc"
	if isGRPCS {
		name = "grpcs"
	}
	cmd := &cobra.Command{
		Use:   name,
		Args:  NotReqArgs,
		Short: fmt.Sprintf("Run server on %s protocal", name),
		Run:   gRPCRun(isGRPCS),
	}
	return cmd
}

func gRPCRun(isGRPCs bool) func(cmd *cobra.Command, args []string) {

	return func(cmd *cobra.Command, args []string) {
		grpc := gRPC(isGRPCs)
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			os.Interrupt,
			syscall.SIGTERM,
			os.Kill)

		switch cc := <-c; cc {
		case os.Interrupt, syscall.SIGTERM, os.Kill:
			log.Info("Signal %+v: interrupt", cc)
		default:
			log.Info("Signal %+v: force stop", cc)
		}
		stopped := make(chan struct{})
		go func() {
			grpc.GracefulStop()
			close(stopped)
		}()
		signal.Notify(c, os.Kill)
		str := ""
		select {
		case <-c:
			grpc.Stop()
			str = "unsafe"
		case <-stopped:
			str = "safe"
		}
		log.Info("Server gRPC is terminated with %s", str)
		return
	}
}

func gRPC(isGRPCs bool) *grpc.Server {
	mod := modules.Get()
	conf := mod.Conf.Svc.App()
	port := conf.Port
	if port == 0 || port > 65535 {
		log.Error("port %d not allow", port)
		os.Exit(1)
	}
	optServer := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 5 * time.Second,
		}),
		grpc.ChainUnaryInterceptor(

			mod.Log.Mid.UnaryServerInterceptor,
		),
		grpc.ChainStreamInterceptor(),
	}
	if isGRPCs {
		pkFile, certFile, okPk, okCert := ssl.CheckSSLPath(conf.SslPrivatePath, conf.SslCertPath)
		if !(okPk && okCert) {
			log.Error("LoadSSLKey not allow")
			os.Exit(1)
		}
		certTLS, err := tls.LoadX509KeyPair(certFile, pkFile)
		if err != nil {
			log.With(log.ErrorString(err)).Error("tls.LoadX509KeyPair")
			os.Exit(1)
		}
		optServer = append(optServer, grpc.Creds(credentials.NewServerTLSFromCert(&certTLS)))
	}

	serve := grpc.NewServer(optServer...)
	routes.RegisterGRPCServer(serve, mod)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.With(log.ErrorString(err)).Error("net.Listen")
		os.Exit(1)
	}
	log.Info("Start gRPC server on :%d", port)
	go func() {
		if err := serve.Serve(lis); err != nil {
			log.With(log.ErrorString(err)).Error("serve.Serve")
			os.Exit(1)
		}
	}()
	return serve
}
