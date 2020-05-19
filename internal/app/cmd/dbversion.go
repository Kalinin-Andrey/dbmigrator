package cmd

import (
	"fmt"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
	"os"

	"github.com/spf13/cobra"
)

// dbversionCmd represents the dbversion command
var dbversionCmd = &cobra.Command{
	Use:   "dbversion",
	Short: "Outputs ID of last applied migration.",
	Long: `Outputs ID of last applied migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dbversion called")
		id, err := dbmigrator.DBVersion()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("dbversion: ", id)
	},
}

func init() {
	rootCmd.AddCommand(dbversionCmd)
}
