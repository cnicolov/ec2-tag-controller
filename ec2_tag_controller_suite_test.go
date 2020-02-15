package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEc2TagController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ec2TagController Suite")
}
