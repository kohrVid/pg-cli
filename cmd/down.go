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

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/kohrVid/pg-cli/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downCmd represents the migrate down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply all down migrations in the database",
	Long:  `Apply all down migrations in the postgres database generated by pg CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		fmt.Println("Using migrations in:", migrationPath)
		conf := viper.Get(env).(map[string]interface{})
		err = db.MigrateDown(conf, migrationPath)
		if err != nil {
			panic(fmt.Errorf("fatal error migrating database: %s \n", err))
		}
	},
}
