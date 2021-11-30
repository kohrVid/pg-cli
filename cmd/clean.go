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

	"github.com/kohrVid/pg-cli/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the postgres database",
	Long:  `This command can be used to delete all rows in the database created by the pg CLI with environment-specific settings`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		conf := viper.Get(env).(map[string]interface{})
		databaseName := conf["database_name"].(string)
		helpers.Clean(conf)

		fmt.Printf("%v database cleaned\n", databaseName)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
