package cmd

import (
	"github.com/buonotti/apisense/v2/filesystem"
	"github.com/buonotti/apisense/v2/log"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize apisense directories",
	Long:  `This command initialize apisense directories. It creates the config directory and the reports and definitions directories.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := filesystem.Setup()
		if err != nil {
			log.DefaultLogger().Fatal(err)
		}
		log.DefaultLogger().Info("Apisense initialized")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
