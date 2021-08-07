module github.com/suhlig/grafana-backup-restore

go 1.16

require (
	github.com/gosimple/slug v1.10.0 // indirect
	github.com/grafana-tools/sdk v0.0.0-20210714133701-11b1efc100c9
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/spf13/cobra v1.2.1
)

replace github.com/grafana-tools/sdk => ../grafana-sdk
