package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator/api"
)

var migrationType	string
var migrationID		uint
var migrationName	string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration.",
	Long: `Create a new migration.`,
	Args: func(cmd *cobra.Command, args []string) error {
		p := &api.MigrationCreateParams{
			ID:		migrationID,
			Type:	migrationType,
			Name:	migrationName,
		}
		return p.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		err := dbmigrator.Create(api.MigrationCreateParams{
			ID:		migrationID,
			Type:	migrationType,
			Name:	migrationName,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().UintVarP(&migrationID, "id", "i", 0, "ID of migration to be created. Must be uint type.")
	createCmd.Flags().StringVarP(&migrationType, "type", "t", "", "Type of migration to be created. Must be one of this: " + fmt.Sprintf("%v", api.MigrationTypes))
	createCmd.Flags().StringVarP(&migrationName, "name", "n", "", "Name of migration to be created. Must be matches the specified regular expression: \"[a-zA-Z0-9_-]+\" ")

	err := rootCmd.MarkFlagRequired("id")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = rootCmd.MarkFlagRequired("type")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = rootCmd.MarkFlagRequired("name")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
