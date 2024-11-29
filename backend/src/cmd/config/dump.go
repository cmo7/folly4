package config

import (
	"fmt"

	"github.com/cmo7/folly4/src/lib/chroma"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump configuration",
	Long:  `Dump configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: folly config dump <file>")
			return
		}
		dumpFile := args[0]
		fmt.Println("Dumping configuration to", dumpFile)
		err := viper.SafeWriteConfigAs(dumpFile)
		if err != nil {
			fmt.Println("Error dumping configuration:", err)
			return
		}

		fmt.Println("Write configuration to", chroma.Color("green")(dumpFile))
	},
}
