/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
)

var migrationType	string
var migrationID		uint
var migrationName	string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		err := sqlmigrator.Create(api.MigrationCreateParams{
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

	rootCmd.MarkFlagRequired("id")
	rootCmd.MarkFlagRequired("type")
	rootCmd.MarkFlagRequired("name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
