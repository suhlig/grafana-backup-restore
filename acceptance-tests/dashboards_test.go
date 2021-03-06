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

var _ = Describe("dashboards", func() {
	var (
		dashboards []map[string]interface{}
	)

	Context("a fresh Grafana server", func() {
		BeforeEach(func() {
			dashboards, err = getDashboards(apiKey)
			Expect(err).ToNot(HaveOccurred())
		})

		It("has no dashboards", func() {
			Expect(dashboards).To(HaveLen(0))
		})
	})

	Context("restoring dashboards", func() {
		JustBeforeEach(func() {
			err = cmd.RestoreDashboards("fixtures/dashboards", "http://localhost:3000", apiKey)
			Expect(err).ToNot(HaveOccurred())
		})

		It("has four dashboards", func() {
			Eventually(func() int {
				dashboards, err = getDashboards(apiKey)
				Expect(err).ToNot(HaveOccurred())
				return len(dashboards)
			}).Should(Equal(4))
		})

		Context("backing up dashboards", func() {
			var targetDirectory string

			BeforeEach(func() {
				targetDirectory, err = ioutil.TempDir("", "grafana-backup")
				Expect(err).ToNot(HaveOccurred())
			})

			JustBeforeEach(func() {
				err = cmd.BackupDashboards(targetDirectory, "http://localhost:3000", apiKey)
				Expect(err).ToNot(HaveOccurred())
			})

			It("has four dashboard files", func() {
				Eventually(func() int {
					count, err := fileCount(targetDirectory)
					Expect(err).ToNot(HaveOccurred())

					return count
				}).Should(Equal(4))
			})

			It("has the expected file names", func() {
				Eventually(func() []string {
					files, err := findFiles(targetDirectory)
					Expect(err).ToNot(HaveOccurred())
					return files
				}).Should(
					SatisfyAll(
						ContainElement(ContainSubstring("General/home.json")),
						ContainElement(ContainSubstring("General/random-data.json")),
						ContainElement(ContainSubstring("Tokyo/nuclear-fallout.json")),
						ContainElement(ContainSubstring("New York/subway-timings.json")),
					),
				)
			})
		})
	})
})

func getDashboards(apiKey string) ([]map[string]interface{}, error) {
	client := &http.Client{Timeout: time.Second * 3}

	req, err := http.NewRequest("GET", "http://localhost:3000/api/search?folderIds=0&query=&starred=false", nil)
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

	var dashboards []map[string]interface{}

	err = json.Unmarshal([]byte(body), &dashboards)

	if err != nil {
		return nil, err
	}

	return dashboards, nil
}
