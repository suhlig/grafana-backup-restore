/*
Copyright © 2021 Steffen Uhlig

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

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:           "grafana-backup-restore",
	Short:         "Simple backup and restore of Grafana",
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// TODO Use Viper instead
	var found bool
	ApiKey, found = os.LookupEnv("GRAFANA_API_TOKEN")

	if !found {
		fmt.Fprintf(os.Stderr, "Error: GRAFANA_API_TOKEN variable not set\n")
		os.Exit(1)
	}

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var Verbose bool
var ApiKey string
var ApiURL string

func init() {
	root.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "verbose output")

	root.AddCommand(backup)
	backup.AddCommand(backupDashboards)
	backup.AddCommand(backupDataSources)
	backup.PersistentFlags().StringVarP(&ApiURL, "url", "U", "", "Grafana API URL (required)")
	backup.MarkPersistentFlagRequired("url")

	root.AddCommand(restore)
	restore.AddCommand(restoreDashboards)
	restore.AddCommand(restoreDataSources)
	restore.PersistentFlags().StringVarP(&ApiURL, "url", "U", "", "Grafana API URL (required)")
	restore.MarkPersistentFlagRequired("url")
}
