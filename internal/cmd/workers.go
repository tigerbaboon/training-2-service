package cmd

import (
	"app/internal/modules/log"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

func closeOsSignal(ch chan os.Signal) {
	defer func() {
		recover()
	}()
	close(ch)
}

// Workers function for call worker
func Workers(work func(chan os.Signal) error, clear func(), timeout time.Duration) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		osSignal := make(chan os.Signal, 1)
		workDone := make(chan os.Signal, 1)
		signal.Notify(osSignal,
			os.Interrupt,
			os.Kill)
		log.Info("Running")
		go func() {
			if err := work(workDone); err != nil {
				log.With(log.ErrorString(err)).Error("Error running")
			}
			close(osSignal)
		}()
		select {
		case s := <-osSignal:
			if s != nil {
				workDone <- s
				switch s {
				case os.Interrupt:
					log.Info("Signal %+v: interrupt", s)
				case os.Kill:
					log.Info("Signal %+v: force stop", s)
				}
				// logger.SetExitCode(128 + int(s.(syscall.Signal)))
			}
			ss := <-osSignal
			if ss != nil {
				log.Info("Exit with SIGKILL")
			} else {
				log.Info("Exit is safe")
			}
			clear()
		}
	}
}
