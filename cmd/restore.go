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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	grafana "github.com/grafana-tools/sdk"
)

func RestoreDataSources(sourceDirectory, apiURL, apiKey string) error {
	return fmt.Errorf("Error: restoring all datasources not yet implemented")
}

func RestoreDashboards(sourceDirectory, apiURL, apiKey string) error {
	absTarget, err := filepath.Abs(sourceDirectory)

	if err != nil {
		return err
	}

	client, err := grafana.NewClient(apiURL, apiKey, grafana.DefaultHTTPClient)

	if err != nil {
		return err
	}

	err = filepath.Walk(absTarget, func(candidate string, info os.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("no info about %s", absTarget)
		}

		if info.IsDir() {
			if Verbose {
				fmt.Fprintf(os.Stderr, "Skipping directory %s\n", candidate)
			}

			return nil
		}

		if filepath.Ext(candidate) != ".json" {
			if Verbose {
				fmt.Fprintf(os.Stderr, "Skipping non-JSON file %s\n", candidate)
			}

			return nil
		}

		relative := path.Dir(candidate[len(absTarget)+1:])

		if Verbose {
			fmt.Fprintf(os.Stderr, "Importing %s into folder %s... ", candidate, relative)
		}

		rawBoard, err := ioutil.ReadFile(candidate)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Skipping %s because it could not be read: %s\n", candidate, err)
			return nil
		}

		ctx := context.Background()
		result, err := client.SetRawDashboard(ctx, rawBoard)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Skipping import of %s: %s\n", candidate, err)
			return nil
		}

		if Verbose {
			fmt.Fprintln(os.Stderr, *result.Status)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
