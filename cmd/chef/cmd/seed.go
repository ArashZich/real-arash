package cmd

import (
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/database"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed data",
	Long:  "seed data to database",
	Run: func(cmd *cobra.Command, args []string) {
		config.Load()
		db := database.Connect()
		database.SeedAllTables(db)
		fmt.Println("Finished!")
	},
}

func init() {
	RootCmd.AddCommand(seedCmd)
}
