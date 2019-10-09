package flaw_test

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw"
)

var _ = Describe("Error", func() {
	It("creates an error successfully", func() {
		err := flaw.New("oh no")
		Expect(err.Error()).To(HavePrefix("message: oh no"))
	})

	It("wraps an error successfully", func() {
		err := flaw.Wrap(fmt.Errorf("oh no"))
		Expect(err.Error()).To(HavePrefix("reason: oh no"))
		Expect(err.Unwrap()).To(MatchError("oh no"))
	})

	Describe("WithCode", func() {
		It("creates an error successfully", func() {
			err := flaw.New("oh no").WithCode(200)
			Expect(err.Error()).To(HavePrefix("code: 200 message: oh no"))
			Expect(err.Code()).To(Equal(200))
		})
	})

	Describe("WithStatus", func() {
		It("creates an error successfully", func() {
			err := flaw.New("oh no").WithStatus(200)
			Expect(err.Status()).To(Equal(200))
		})
	})

	Describe("WithError", func() {
		It("creates an error successfully", func() {
			err := flaw.New("failed").WithError(fmt.Errorf("oh no"))
			Expect(err.Error()).To(HavePrefix("message: failed reason: oh no"))
			Expect(err.Unwrap()).To(MatchError("oh no"))
		})
	})

	Describe("WithMessage", func() {
		It("creates an error successfully", func() {
			err := flaw.Wrap(fmt.Errorf("oh no")).WithMessage("failed")
			Expect(err.Error()).To(HavePrefix("message: failed reason: oh no"))
			Expect(err.Unwrap()).To(MatchError("oh no"))
		})
	})

	Describe("Format", func() {
		It("prints the error successfully", func() {
			err := flaw.New("failed").WithCode(404).WithError(fmt.Errorf("oh no"))
			Expect(fmt.Sprintf("%v", err)).To(HavePrefix("code: 404 message: failed reason: oh no stack:"))
		})

		Context("when the verbose printing is used", func() {
			It("prints the error successfully", func() {
				err := flaw.New("failed").WithCode(404).WithError(fmt.Errorf("oh no"))
				Expect(fmt.Sprintf("%+v", err)).To(HavePrefix("code: 404\nmessage: failed\nreason: oh no\nstack:"))
			})
		})
	})

	Describe("MarshalJSON", func() {
		It("marshals the error successfully", func() {
			errx := flaw.New("oh no").WithCode(200)
			errx.Wrap(fmt.Errorf("failed"))

			data, err := json.Marshal(errx)
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(`{"code":200,"message":"oh no","reason":"failed"}`))
		})

		Context("when the wrapped error implements MarshalJSON", func() {
			It("marshals the error successfully", func() {
				errx := flaw.New("oh no").WithCode(200)
				errx.Wrap(flaw.New("failed"))

				data, err := json.Marshal(errx)
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`{"code":200,"message":"oh no","reason":{"message":"failed"}}`))
			})
		})
	})
})

var _ = Describe("ErrorCollection", func() {
	It("creates an error successfully", func() {
		errs := flaw.ErrorCollector{}
		errs = append(errs, fmt.Errorf("oh no"))
		Expect(errs).To(MatchError("oh no"))
	})

	Context("when the collector has more than one error", func() {
		It("creates an error successfully", func() {
			errs := flaw.ErrorCollector{}
			errs = append(errs, fmt.Errorf("oh no"))
			errs = append(errs, fmt.Errorf("oh yes"))
			Expect(errs).To(MatchError("oh no; oh yes; "))
		})
	})

	Describe("Unwrap", func() {
		It("unwraps the first error", func() {
			errs := flaw.ErrorCollector{}
			errs = append(errs, fmt.Errorf("oh no"))
			Expect(errs.Unwrap()).To(MatchError("oh no"))
		})

		Context("when the collector is empty", func() {
			It("unwraps the nil error", func() {
				errs := flaw.ErrorCollector{}
				Expect(errs.Unwrap()).To(BeNil())
			})
		})

		Context("when the collector has more than one error", func() {
			Describe("Unwrap", func() {
				It("unwraps the errors as itself", func() {
					errs := flaw.ErrorCollector{}
					errs = append(errs, fmt.Errorf("oh no"))
					errs = append(errs, fmt.Errorf("oh yes"))
					Expect(errs.Unwrap()).To(Equal(errs))
				})
			})
		})
	})

	Describe("Format", func() {
		It("prints the error successfully", func() {
			errs := flaw.ErrorCollector{}
			errs = append(errs, fmt.Errorf("oh no"))
			errs = append(errs, fmt.Errorf("oh yes"))
			Expect(fmt.Sprintf("%v", errs)).To(Equal("oh no; oh yes; "))
		})

		Context("when the verbose printing is used", func() {
			It("prints the error successfully", func() {
				errs := flaw.ErrorCollector{}
				errs = append(errs, fmt.Errorf("oh no"))
				errs = append(errs, fmt.Errorf("oh yes"))
				Expect(fmt.Sprintf("%+v", errs)).To(Equal(" --- oh no\n --- oh yes\n"))
			})
		})
	})

	Describe("MarshalJSON", func() {
		It("marshals the error successfully", func() {
			errs := flaw.ErrorCollector{}
			errs = append(errs, fmt.Errorf("oh no"))

			data, err := json.Marshal(errs)
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(`["oh no"]`))
		})

		Context("when the child error implements MarshalJSON", func() {
			It("marshals the error successfully", func() {
				errs := flaw.ErrorCollector{}
				errs = append(errs, flaw.New("oh no"))

				data, err := json.Marshal(errs)
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`[{"message":"oh no"}]`))
			})
		})
	})
})
