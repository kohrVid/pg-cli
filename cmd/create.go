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
	"log"

	"github.com/kohrVid/pg-cli/helpers"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new database for the calendar API",
	Long:  `This command can be used to create a database with environment-specific settings`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		conf := viper.Get(env).(map[string]interface{})
		databaseUser := conf["database_user"].(string)
		databaseName := conf["database_name"].(string)

		createRole := fmt.Sprintf("CREATE ROLE %v", databaseUser)
		alterRole := fmt.Sprintf(
			"ALTER ROLE %v WITH SUPERUSER LOGIN CREATEDB;",
			databaseUser,
		)

		createDB := fmt.Sprintf(
			"CREATE DATABASE %v WITH OWNER %v ENCODING 'UTF8';",
			databaseName,
			databaseUser,
		)

		db := helpers.PostgresDB()
		defer db.Close()

		_, err = db.Exec(createRole)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec(alterRole)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(createDB)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%v database created\n", databaseName)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
