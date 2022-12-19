package flaw_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFlaw(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flaw Suite")
}
