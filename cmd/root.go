/*
Copyright © 2021 Jessica Été <kohrVid@zoho.com>

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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var env string
var migrationPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pg-cli",
	Short: "Database Operations",
	Long:  `This tool can be used to create and manage a postgres database`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config",
		"c", "",
		"config file (default is $GOPATH/src/github.com/example/config/env.yaml)",
	)

	rootCmd.MarkFlagRequired("config")

	createCmd.Flags().StringVarP(&env, "env", "e", "development", "environment name")
	cleanCmd.Flags().StringVarP(&env, "env", "e", "development", "environment name")
	dropCmd.Flags().StringVarP(&env, "env", "e", "development", "environment name")
	seedCmd.Flags().StringVarP(&env, "env", "e", "development", "environment name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		examplePath := filepath.Join(
			os.Getenv("GOPATH"), "src",
			"github.com", "example",
			"config",
		)

		viper.AddConfigPath(examplePath)
		viper.SetConfigName("env")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
