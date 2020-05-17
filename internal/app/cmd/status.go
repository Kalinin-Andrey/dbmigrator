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
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")
		ms, err := sqlmigrator.Status()
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
