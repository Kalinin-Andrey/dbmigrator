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
	"context"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/sqlmigrator/api"

	_ "github.com/Kalinin-Andrey/dbmigrator/migration"
)

var cfgFile, logFile, dsn, dir string
var ctx context.Context

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbmigrator",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(c context.Context) {
	ctx = c

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initSQLMigrator)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dbmigrator.yaml)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log", "", "log file (default is stdout)")
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "dsn string for connection to DB")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", "", "path to directory with migrations")

	err := viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = viper.BindPFlag("dsn", rootCmd.PersistentFlags().Lookup("dsn"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// initSQLMigrator reads in config file and ENV variables if set.
func initSQLMigrator() {
	var config api.Configuration

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dbmigrator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dbmigrator")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		fmt.Println("Default config file not found:", viper.ConfigFileUsed())
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("\n Error: expected arguments \"dsn\" and \"dir\" was not set\n", rootCmd.Help())
		os.Exit(1)
	}
	config.ExpandEnv()

	err = sqlmigrator.Init(ctx, config, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

