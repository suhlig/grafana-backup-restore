/*
Copyright 2016 Alexander I.Grafov <grafov@gmail.com>
Copyright 2016-2019 The Grafana SDK authors
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
	"context"
	"fmt"
	"io/ioutil"
	"os"

	grafana "github.com/grafana-tools/sdk"
	"github.com/spf13/cobra"
)

var backup = &cobra.Command{
	Use:   "backup",
	Short: "Backup Grafana items",
}

var backupDashboards = &cobra.Command{
	Use:           "dashboards",
	Short:         "Backup all dashboards",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			dashboards []grafana.FoundBoard
			data       []byte
			meta       grafana.BoardProperties
			err        error
		)

		c, err := grafana.NewClient(ApiURL, ApiKey, grafana.DefaultHTTPClient)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
			os.Exit(1)
		}

		ctx := context.Background()

		if dashboards, err = c.SearchDashboards(ctx, "", false); err != nil {
			return err
		}

		for _, dashboard := range dashboards {
			data, meta, err = c.GetRawDashboardByUID(ctx, dashboard.UID)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s for %s\n", err, dashboard.URI)
				continue
			}

			var directory, displayPath string

			if dashboard.FolderTitle == "" {
				directory = meta.FolderTitle
				displayPath = fmt.Sprintf("%s/%s", meta.FolderTitle, dashboard.Title)
			} else {
				directory = dashboard.FolderTitle
				displayPath = fmt.Sprintf("%s/%s", dashboard.FolderTitle, dashboard.Title)
			}

			err = os.MkdirAll(directory, os.FileMode(int(0700)))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating backup folder %s: %s\n", directory, err)
			}

			fileName := fmt.Sprintf("%s/%s.json", directory, meta.Slug)

			if Verbose {
				fmt.Fprintf(os.Stderr, "Writing dashboard '%s' to %s\n", displayPath, fileName)
			}

			err = ioutil.WriteFile(fileName, data, os.FileMode(int(0600)))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing %s: %s\n", fileName, err)
			}
		}

		return nil
	},
}

var backupDataSources = &cobra.Command{
	Use:           "datasources",
	Short:         "Backup all datasources",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("Error: backing up all datasources not yet implemented")
	},
}
