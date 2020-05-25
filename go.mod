module github.com/bigkevmcd/tekton-archiver

go 1.14

require (
	github.com/aws/aws-sdk-go v1.27.1
	github.com/google/go-cmp v0.4.0
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/tektoncd/pipeline v0.12.1
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975 // indirect
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c // indirect
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89 // indirect
	knative.dev/pkg v0.0.0-20200306230727-a56a6ea3fa56
	sigs.k8s.io/yaml v1.2.0 // indirect
)

// Knative deps (release-0.13)
replace (
	contrib.go.opencensus.io/exporter/stackdriver => contrib.go.opencensus.io/exporter/stackdriver v0.12.9-0.20191108183826-59d068f8d8ff
	knative.dev/caching => knative.dev/caching v0.0.0-20200116200605-67bca2c83dfa
	knative.dev/pkg => knative.dev/pkg v0.0.0-20200306230727-a56a6ea3fa56
	knative.dev/pkg/vendor/github.com/spf13/pflag => github.com/spf13/pflag v1.0.5
)

// Pin k8s deps to 1.16.5
replace (
	k8s.io/api => k8s.io/api v0.16.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.5
	k8s.io/client-go => k8s.io/client-go v0.16.5
	k8s.io/code-generator => k8s.io/code-generator v0.16.5
	k8s.io/gengo => k8s.io/gengo v0.0.0-20190327210449-e17681d19d3a
)
