package flaw_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw"
)

var _ = Describe("StackTrace", func() {
	Describe("Format", func() {
		It("prints the stack trace successfully", func() {
			stack := flaw.NewStackTrace()
			Expect(fmt.Sprintf("%v", stack)).To(ContainSubstring("[runner.go:113"))
		})

		It("prints the stack trace verbosely", func() {
			stack := flaw.NewStackTrace()
			Expect(fmt.Sprintf("%+v", stack)).To(ContainSubstring("github.com/onsi/ginkgo@v1.10.2/internal/leafnodes/runner.go"))
		})

		It("prints the stack trace shortly", func() {
			stack := flaw.NewStackTrace()
			Expect(fmt.Sprintf("%s", stack)).To(HavePrefix("[runner.go"))
		})

		It("prints the stack trace verbosely as slice", func() {
			stack := flaw.NewStackTrace()
			Expect(fmt.Sprintf("%+s", stack)).To(ContainSubstring("github.com/onsi/ginkgo@v1.10.2/internal/leafnodes/runner.go"))
		})
	})
})
