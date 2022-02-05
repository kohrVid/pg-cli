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

// dropCmd represents the drop command
var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop the database",
	Long:  `This command can be used to drop the database created by the pg CLI with environment-specific settings`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		conf := viper.Get(env).(map[string]interface{})
		databaseUser := conf["database_user"].(string)
		databaseName := conf["database_name"].(string)

		dropRole := fmt.Sprintf("DROP ROLE %v", databaseUser)
		dropDB := fmt.Sprintf("DROP DATABASE %v;", databaseName)

		conn := db.PgConnect()
		defer conn.Close()

		_, err = conn.Exec(dropDB)
		if err != nil {
			fmt.Println(err)
		}

		_, err = conn.Exec(dropRole)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%v database deleted\n", databaseName)
	},
}

func init() {
	rootCmd.AddCommand(dropCmd)
}
