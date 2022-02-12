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

	"github.com/kohrVid/pg-cli/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "See the postgres database",
	Long:  `This command can be used to seed a database created by the pg CLI with environment-specific settings`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		conf := viper.Get(env).(map[string]interface{})

		err = db.Seed(conf)
		if err != nil {
			panic(fmt.Errorf("Fatal error seeding database: %s \n", err))
		}

		fmt.Printf(
			"Seeded %v database\n",
			conf["database_name"].(string),
		)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
