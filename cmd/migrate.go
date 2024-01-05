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
	"github.com/spf13/cobra"
)

var migrationPath string
var step int

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the postgres database",
	Long:  `This command can be used to migrate the database created by the pg CLI with environment-specific settings`,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.AddCommand(stepCmd)
	stepCmd.Flags().StringVarP(&migrationPath, "path", "p", "./migrations", "relative path to the migration files")
	stepCmd.Flags().IntVarP(&step, "", "n", 0, "number of steps to migrate")
	stepCmd.MarkFlagRequired("step")

	migrateCmd.AddCommand(migrateVersionCmd)
	migrateVersionCmd.Flags().StringVarP(&migrationPath, "path", "p", "./migrations", "relative path to the migration files")

	migrateCmd.AddCommand(upCmd)
	upCmd.Flags().StringVarP(&migrationPath, "path", "p", "./migrations", "relative path to the migration files")

	migrateCmd.AddCommand(downCmd)
	downCmd.Flags().StringVarP(&migrationPath, "path", "p", "./migrations", "relative path to the migration files")
}
