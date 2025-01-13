package main

import (
	"app/app/console"
	"app/app/modules"
	"app/internal/cmd"
	"app/internal/service/provider"
	"log/slog"
	"os"

	_ "time/tzdata"

	"github.com/spf13/cobra"
	_ "google.golang.org/grpc/encoding/gzip"
	// _ "google.golang.org/grpc/xds"
)

func main() {
	cobra.EnableCommandSorting = false
	if err := exec(); err != nil {
		slog.Error("Error running")
		os.Exit(1)
	}
}

func command() error {
	cmda := &cobra.Command{
		Use: "app",
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			pv := provider.Config(modules.Map())
			return pv.Close(cmd.Context())
		},
		Args: cmd.NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmds := &cobra.Command{
		Use:   "cmd",
		Short: "List all commands",
	}
	cmds.AddCommand(console.Commands()...)

	cmda.AddCommand(cmd.HTTP(false), cmd.HTTP(true))
	cmda.AddCommand(cmd.GRPC(false), cmd.GRPC(true))
	cmda.AddCommand(cmd.Migrate())
	cmda.AddCommand(cmd.Module())
	cmda.AddCommand(cmds)
	return cmda.Execute()
}
