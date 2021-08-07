package acceptance_tests_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var apiKey string
var err error

var _ = BeforeSuite(func() {
	apiKey, err = createApiKey()
	Expect(err).ToNot(HaveOccurred())
})

func TestAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Tests Suite")
}

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

func fileCount(directory string) (int, error) {
	fileCount := 0

	err := filepath.Walk(directory, func(candidate string, info os.FileInfo, err error) error {
		if !info.IsDir() || filepath.Ext(candidate) == ".json" {
			fileCount += 1
		}
		return nil
	})

	return fileCount, err
}

func findFiles(directory string) ([]string, error) {
	var fileNames []string

	err := filepath.Walk(directory, func(candidate string, info os.FileInfo, err error) error {
		if !info.IsDir() || filepath.Ext(candidate) == ".json" {
			fileNames = append(fileNames, candidate)
		}
		return nil
	})

	return fileNames, err
}
