module github.com/cnicolov/ec2-tag-controller

go 1.13

require (
	github.com/aws/aws-sdk-go v1.29.3
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.0.0-20200127113903-12be8a0d907a
	sigs.k8s.io/controller-runtime v0.5.0
)
