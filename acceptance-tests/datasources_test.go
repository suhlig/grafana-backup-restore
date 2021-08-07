package acceptance_tests_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/suhlig/grafana-backup-restore/cmd"
)

var _ = Describe("datasources", func() {
	var (
		datasources []map[string]interface{}
	)

	Context("a fresh Grafana server", func() {
		BeforeEach(func() {
			datasources, err = getDatasources(apiKey)
			Expect(err).ToNot(HaveOccurred())
		})

		It("has no dashboards", func() {
			Expect(datasources).To(HaveLen(0))
		})
	})

	Context("restoring datasources", func() {
		JustBeforeEach(func() {
			err = cmd.RestoreDatasources("fixtures/datasources", "http://localhost:3000", apiKey)
			Expect(err).ToNot(HaveOccurred())
		})

		It("has two datasources", func() {
			Eventually(func() int {
				datasources, err = getDatasources(apiKey)
				Expect(err).ToNot(HaveOccurred())
				return len(datasources)
			}).Should(Equal(2))
		})

		Context("backing up datasources", func() {
			var targetDirectory string

			BeforeEach(func() {
				targetDirectory, err = ioutil.TempDir("", "grafana-backup")
				Expect(err).ToNot(HaveOccurred())
			})

			JustBeforeEach(func() {
				err = cmd.BackupDatasources(targetDirectory, "http://localhost:3000", apiKey)
				Expect(err).ToNot(HaveOccurred())
			})

			It("has two datasource files", func() {
				Eventually(func() int {
					count, err := fileCount(targetDirectory)
					Expect(err).ToNot(HaveOccurred())

					return count
				}).Should(Equal(2))
			})

			It("has the expected file names", func() {
				Eventually(func() []string {
					files, err := findFiles(targetDirectory)
					Expect(err).ToNot(HaveOccurred())
					return files
				}).Should(
					SatisfyAll(
						ContainElement(ContainSubstring("TestData DB 0.json")),
						ContainElement(ContainSubstring("TestData DB 1.json")),
					),
				)
			})
		})
	})
})

func getDatasources(apiKey string) ([]map[string]interface{}, error) {
	client := &http.Client{Timeout: time.Second * 3}

	req, err := http.NewRequest("GET", "http://localhost:3000/api/datasources", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status was not 200, but %d", resp.StatusCode)
	}

	var datasources []map[string]interface{}

	err = json.Unmarshal([]byte(body), &datasources)

	if err != nil {
		return nil, err
	}

	return datasources, nil
}
