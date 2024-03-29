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

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stepCmd represents the migrate step command
var stepCmd = &cobra.Command{
	Use:   "step -n <int>",
	Short: "Migrate a database in steps",
	Long: `This command can be used to migrate a database created by the pg CLI by a given number of steps.

Steps can be positive or negative so as to indicate the direction of the
migration (e.g., ` + "`" + `step -n 1` + "`" + ` would go up a single migration whereas ` + "`" + `step -n -1` + "`" + `
would go down a single migration).`,
	//+ "`" + `step -n 1` + "`" + ` would go
	//up a single migration wheras `step -n -1` would go down a single migration).
	//`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		fmt.Println("Using migrations in:", migrationPath)
		conf := viper.Get(env).(map[string]interface{})
		err = db.MigrateStep(conf, migrationPath, step)
		if err != nil {
			panic(fmt.Errorf("Fatal error migrating database: %s \n", err))
		}
	},
}
