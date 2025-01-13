package cmd

import (
	"app/internal/service/http"
	"fmt"

	"github.com/spf13/cobra"
)

// HTTP is serve http ot https
func HTTP(isHTTPS bool) *cobra.Command {
	name := "http"
	if isHTTPS {
		name = "https"
	}
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Run server on %s protocal", name),
		Run:   http.HTTPD(isHTTPS),
	}
	return cmd
}
