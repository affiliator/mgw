// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/affiliator/mgw/storage"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database related tools. Use ./mgw db -h for more.",
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database schema. Can be used to migrate skeleton.",
	Long: `WARNING: AutoMigrate will ONLY create tables, missing columns and missing indexes,
and WON’T change existing column’s type or delete unused columns to protect your data.`,
	Run: migrate,
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbMigrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	storage.Connection().AutoMigrate(&storage.Message{})
}
