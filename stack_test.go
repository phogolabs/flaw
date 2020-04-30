package flaw_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw"
)

var _ = Describe("StackTrace", func() {
	Describe("NewStackTraceAt", func() {
		It("skips a frame successfully", func() {
			var (
				stack   = flaw.NewStackTrace()
				skipped = flaw.NewStackTraceAt(1)
			)

			Expect(len(stack)).To(Equal(len(skipped) + 1))
		})
	})

	Describe("Format", func() {
		It("prints the stack successfully", func() {
			stack := flaw.NewStackTrace()
			Expect(fmt.Sprintf("%v", stack)).To(ContainSubstring("leafnodes/runner.go"))
		})

		Context("when the %#v format is used", func() {
			It("prints the stack successfully", func() {
				stack := flaw.NewStackTrace()
				Expect(fmt.Sprintf("%#v", stack)).To(ContainSubstring("[]flaw.StackFrame{"))
			})
		})

		Context("when the %+v format is used", func() {
			It("prints the stack successfully", func() {
				stack := flaw.NewStackTrace()
				Expect(fmt.Sprintf("%+v", stack)).To(ContainSubstring("leafnodes/runner.go"))
			})
		})

		Context("when the %s format is used", func() {
			It("prints the stack successfully", func() {
				stack := flaw.NewStackTrace()
				Expect(fmt.Sprintf("%s", stack)).To(HavePrefix("[leafnodes/runner.go"))
			})
		})

		Context("when the %+s format is used", func() {
			It("prints the stack successfully", func() {
				stack := flaw.NewStackTrace()
				Expect(fmt.Sprintf("%+s", stack)).To(ContainSubstring("leafnodes/runner.go"))
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
			Expect(string(data)).To(HavePrefix("leafnodes/runner.go"))
		})
	})
})
