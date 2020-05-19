package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
)

// redoCmd represents the redo command
var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Starts down and then up actions of one last migration.",
	Long: `Starts down and then up actions of one last migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("redo called")
		err := dbmigrator.Redo()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(redoCmd)

}
