package cmd

import (
	"context"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator"
	"github.com/Kalinin-Andrey/dbmigrator/pkg/dbmigrator/api"

	//_ "github.com/Kalinin-Andrey/dbmigrator/migration"
)

var cfgFile, logFile, dsn, dir string
var ctx context.Context

var rootCmd = &cobra.Command{
	Use:   "dbmigrator",
	Short: "Go database migration tool and library supporting SQL or Go scripts.",
	Long: `Go database migration tool and library supporting SQL or Go scripts.`,
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

	//err = dbmigrator.Init(ctx, config, nil)
	err = dbmigrator.InitTool(ctx, config, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

