package cmd

import (
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/database"
	"github.com/spf13/cobra"
)

var migrateDown = &cobra.Command{
	Use:   "migrate:down",
	Short: "migrate down",
	Long:  "migrate down",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.Connect()
		database.DropAll(db)
		fmt.Println("Finished!")
	},
}

func init() {
	RootCmd.AddCommand(migrateDown)
}
