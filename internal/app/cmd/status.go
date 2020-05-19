package cmd

import (
	"fmt"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator/api"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Outputs saved in db logs of migrations.",
	Long: `Outputs saved in db logs of migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")
		ms, err := dbmigrator.Status()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printHeader()

		for _, m := range ms {
			fmt.Printf("| %6d | %-50s | %11s | %v |\n", m.ID, m.Name, api.MigrationStatuses[int(m.Status)], m.Time)
		}

		printLine()
	},
}

func printHeader() {
	fmt.Println("Status of migrations")
	printLine()
	fmt.Printf("| %6s | %-50s | %11s | %-36s |\n", "ID", "Name", "Status", "Time")
	printLine()
}

func printLine() {
	fmt.Println(strings.Repeat("-", 116))
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
