package console

import (
	"app/app/modules"
	"app/internal/cmd"
	"app/internal/modules/log"

	"github.com/spf13/cobra"
)

func helloCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "hello",
		Args: cmd.NotReqArgs,
		Run: func(cmd *cobra.Command, args []string) {
			modules := modules.Get()
			log.Info("Hello, world")
			for k, v := range modules.Map() {
				log.Info("%s %#v", k, v)
			}
		},
	}
	return cmd
}
