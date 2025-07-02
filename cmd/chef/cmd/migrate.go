package cmd

import (
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/database"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	Long:  "migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.Connect()
		database.Migrate(db)
		fmt.Println("Migration finished!")
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
