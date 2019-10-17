package format_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw/format"
)

var _ = Describe("StringSlice", func() {
	slice := format.StringSlice{"hello", "world"}

	It("formats the slice successfully", func() {
		Expect(fmt.Sprintf("%v", slice)).To(Equal("[hello, world]"))
	})

	Context("when the format %s is used", func() {
		It("formats the slice successfully", func() {
			Expect(fmt.Sprintf("%s", slice)).To(Equal("[hello, world]"))
		})
	})

	Context("when the format %+v is used", func() {
		It("formats the slice successfully", func() {
			Expect(fmt.Sprintf("%+v", slice)).To(Equal(" --- hello\n --- world"))
		})
	})

	Context("when the format %#v is used", func() {
		It("formats the slice successfully", func() {
			Expect(fmt.Sprintf("%#v", slice)).To(Equal("[]string{\"hello\", \"world\"}"))
		})
	})
})
