package acceptance_tests_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backup and restore dashboards", func() {
	var (
		err        error
		apiKey     string
		dashboards []map[string]interface{}
	)

	BeforeSuite(func() {
		apiKey, err = createApiKey()
		Expect(err).ToNot(HaveOccurred())
	})

	Context("a fresh Grafana server", func() {
		BeforeEach(func() {
			dashboards, err = getDashboards(apiKey)
			Expect(err).ToNot(HaveOccurred())
		})

		It("has no dashboards", func() {
			Expect(dashboards).To(HaveLen(0))
		})
	})
})

func createApiKey() (string, error) {
	response, err := http.PostForm("http://admin:admin@localhost:3000/api/auth/keys", url.Values{"name": {fmt.Sprintf("Acceptance Test %d", GinkgoRandomSeed())}, "role": {"Admin"}})

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)

	if err != nil {
		return "", err
	}

	return result["key"].(string), nil
}

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
