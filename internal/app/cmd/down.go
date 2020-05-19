package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Starts down action of one last migration.",
	Long: `Starts down action of one last migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("down called")
		err := dbmigrator.Down(0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

}
