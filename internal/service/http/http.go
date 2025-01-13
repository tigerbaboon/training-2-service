package http

import (
	"app/app/modules"
	"app/internal/modules/log"
	"app/routes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ginLogger struct{}

func (l *ginLogger) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return len(p), nil
}

// HTTPD is the main function for the HTTP server.
func HTTPD(isHTTPS bool) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx, cancel := NotifyContext()
		mod := modules.Get()
		conf := mod.Conf.Svc.App()

		srv := serve(ctx, mod)

		go func() {
			srv.Addr = fmt.Sprintf("0.0.0.0:%d", conf.Port)
			log.Info("Gin is running on %s.", srv.Addr)
			if isHTTPS {
				pk := conf.SslPrivatePath
				cert := conf.SslCertPath
				if err := srv.ListenAndServeTLS(cert, pk); !errors.Is(err, nil) && !errors.Is(err, http.ErrServerClosed) {
					log.With(log.ErrorString(err)).Error("Gin was failed to start %s.", srv.Addr)
					os.Exit(1)
				}
			} else if err := srv.ListenAndServe(); !errors.Is(err, nil) && !errors.Is(err, http.ErrServerClosed) {
				log.With(log.ErrorString(err)).Error("Gin was failed to start %s.", srv.Addr)
				os.Exit(1)
			}
		}()

		<-ctx.Done()
		cancel()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(timeoutCtx); err != nil {
			log.With(log.ErrorString(err)).Error("Gin was failed to shutdown.")
			os.Exit(1)
		}
		log.Info("Gin was successful shutdown.")
	}
}

func serve(ctx context.Context, mod *modules.Modules) *http.Server {
	conf := mod.Conf.Svc.App()
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()

	if conf.Debug {
		app.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Output: new(ginLogger),
			Formatter: func(params gin.LogFormatterParams) string {
				if params.ErrorMessage != "" {
					return params.ErrorMessage
				}
				return fmt.Sprintf("%d %s %s %s %s %s %s",
					params.StatusCode,
					params.Method,
					params.Path,
					params.Latency,
					params.ClientIP,
					params.Request.Proto,
					params.Request.UserAgent(),
				)
			},
		}), gin.Recovery())
		pprof.Register(app)
	}
	routes.Router(app, mod)

	h2s := &http2.Server{}

	srv := &http.Server{
		Handler: h2c.NewHandler(app, h2s),
	}
	srv.SetKeepAlivesEnabled(true)

	srv.RegisterOnShutdown(func() {

	})
	return srv
}
