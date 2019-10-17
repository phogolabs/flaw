package flaw_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw"
)

var _ = Describe("StackTrace", func() {
	Describe("Skip", func() {
		It("skips a frame successfully", func() {
			stack := flaw.NewStackTrace()
			count := len(stack)
			stack.Skip(1)
			Expect(count).To(Equal(len(stack) + 1))
		})
	})

	Describe("Format", func() {
		It("prints the stack successfully", func() {
			stack := flaw.NewStackTrace()
			Expect(fmt.Sprintf("%v", stack)).To(ContainSubstring("[runner.go:113"))
		})

		Context("when the %#v format is used", func() {
			It("prints the stack successfully", func() {
				stack := flaw.NewStackTrace()
				Expect(fmt.Sprintf("%#v", stack)).To(ContainSubstring("[]flaw.StackFrame{runner.go:113"))
			})
		})

		Context("when the %+v format is used", func() {
			It("prints the stack successfully", func() {
				stack := flaw.NewStackTrace()
				Expect(fmt.Sprintf("%+v", stack)).To(ContainSubstring("github.com/onsi/ginkgo@v1.10.2/internal/leafnodes/runner.go"))
			})
		})

		Context("when the %s format is used", func() {
			It("prints the stack successfully", func() {
				stack := flaw.NewStackTrace()
				Expect(fmt.Sprintf("%s", stack)).To(HavePrefix("[runner.go"))
			})
		})

		Context("when the %+s format is used", func() {
			It("prints the stack successfully", func() {
				stack := flaw.NewStackTrace()
				Expect(fmt.Sprintf("%+s", stack)).To(ContainSubstring("github.com/onsi/ginkgo@v1.10.2/internal/leafnodes/runner.go"))
			})
		})
	})
})

var _ = Describe("StackFrame", func() {
	Describe("MarshalText", func() {
		It("prints the frame successfully", func() {
			frame := flaw.NewStackTrace()[0]
			data, err := frame.MarshalText()
			Expect(err).To(BeNil())
			Expect(string(data)).To(HaveSuffix("github.com/onsi/ginkgo@v1.10.2/internal/leafnodes/runner.go:113 (github.com/onsi/ginkgo/internal/leafnodes.(*runner).runSync)"))
		})
	})
})
