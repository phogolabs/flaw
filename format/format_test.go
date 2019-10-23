package format_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw/format"
	"github.com/phogolabs/flaw/format/fake"
)

var _ = Describe("State", func() {
	var (
		state   *format.State
		flusher *fake.StateFlusher
	)

	BeforeEach(func() {
		flusher = &fake.StateFlusher{}
	})

	JustBeforeEach(func() {
		state = format.NewState(flusher)
	})

	Context("when the flusher has + flag", func() {
		BeforeEach(func() {
			flusher.FlagReturns(true)
		})

		It("creates a tab writer", func() {
			_, err := state.Write([]byte("hello"))
			Expect(err).To(BeNil())
			Expect(state.Size()).To(Equal(5))
		})
	})

	Describe("Size", func() {
		BeforeEach(func() {
			flusher.WriteReturns(10, nil)
		})

		It("returns the value successfully", func() {
			_, err := state.Write([]byte("hello"))
			Expect(err).To(BeNil())
			Expect(state.Size()).To(Equal(10))
		})
	})

	Describe("Width", func() {
		BeforeEach(func() {
			flusher.WidthReturns(10, true)
		})

		It("returns the value successfully", func() {
			size, ok := state.Width()
			Expect(ok).To(BeTrue())
			Expect(size).To(Equal(10))
		})
	})

	Describe("Precision", func() {
		BeforeEach(func() {
			flusher.PrecisionReturns(10, true)
		})

		It("returns the value successfully", func() {
			size, ok := state.Precision()
			Expect(ok).To(BeTrue())
			Expect(size).To(Equal(10))
		})
	})

	Describe("Flag", func() {
		BeforeEach(func() {
			flusher.FlagReturns(true)
		})

		It("returns the flag successfully", func() {
			Expect(state.Flag('+')).To(BeTrue())
		})
	})

	Describe("Write", func() {
		It("writes the content successfully", func() {
			_, err := state.Write([]byte("hello"))
			Expect(err).To(BeNil())
			Expect(flusher.WriteCallCount()).To(Equal(1))
			Expect(flusher.WriteArgsForCall(0)).To(Equal([]byte("hello")))
		})
	})

	Describe("Flush", func() {
		It("flushes the writer successfully", func() {
			Expect(state.Flush()).To(Succeed())
			Expect(flusher.FlushCallCount()).To(Equal(1))
		})

		Context("when the writer is not flusher", func() {
			JustBeforeEach(func() {
				state = format.NewState(nil)
			})

			It("flushes the writer successfully", func() {
				Expect(state.Flush()).To(Succeed())
			})
		})
	})
})

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
