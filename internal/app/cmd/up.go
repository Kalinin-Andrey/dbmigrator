package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"

)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Starts up actions of migrations.",
	Long: `Starts up actions of migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("up called")
		err := dbmigrator.Up(0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

}
