package cmd

import (
	"github.com/buonotti/apisense/v2/api/db"
	"github.com/buonotti/apisense/v2/log"
	"github.com/spf13/cobra"
)

var apiUserDisableCmd = &cobra.Command{
	Use:               "disable [USERNAME]",
	Short:             "Disable a user",
	Long:              `This command allows to disable a user of the API.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validEnabledUserFunc(),
	Run: func(_ *cobra.Command, args []string) {
		username := args[0]

		err := db.DisableUser(username)
		if err != nil {
			log.DefaultLogger().Fatal(err)
		}

		log.DefaultLogger().Info("User disabled", "username", username)
	},
}

func init() {
	apiUserCmd.AddCommand(apiUserDisableCmd)
}
