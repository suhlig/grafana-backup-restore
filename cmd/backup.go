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

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup Grafana items",
}

var backupDashboardsCmd = &cobra.Command{
	Use:           "dashboards",
	Short:         "Backup all dashboards",
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			boardLinks []grafana.FoundBoard
			rawBoard   []byte
			meta       grafana.BoardProperties
			err        error
		)

		c, err := grafana.NewClient(ApiURL, ApiKey, grafana.DefaultHTTPClient)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
			os.Exit(1)
		}

		ctx := context.Background()

		if boardLinks, err = c.SearchDashboards(ctx, "", false); err != nil {
			return err
		}

		for _, link := range boardLinks {
			rawBoard, meta, err = c.GetRawDashboardByUID(ctx, link.UID)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s for %s\n", err, link.URI)
				continue
			}

			err = os.MkdirAll(meta.FolderTitle, os.FileMode(int(0700)))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating backup folder %s: %s\n", meta.FolderTitle, err)
			}

			fileName := fmt.Sprintf("%s/%s.json", meta.FolderTitle, meta.Slug)

			if Verbose {
				fmt.Fprintf(os.Stderr, "Writing dashboard '%s/%s' to %s\n", link.FolderTitle, link.Title, fileName)
			}

			err = ioutil.WriteFile(fileName, rawBoard, os.FileMode(int(0600)))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing %s: %s\n", fileName, err)
			}
		}

		return nil
	},
}

var backupDataSourceCmd = &cobra.Command{
	Use:          "datasource",
	Short:        "Backup a datasource",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("Error: backing up a datasource not yet implemented")
	},
}
