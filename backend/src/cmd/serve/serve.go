package serve

import (
	"github.com/cmo7/folly4/src/app"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long:  `Start the web server`,
	Run: func(cmd *cobra.Command, args []string) {
		app.Serve()
	},
}

/**
*
* Te amo :3
*
 */
