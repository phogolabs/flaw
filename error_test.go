package flaw_test

import (
	"encoding/json"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw"
)

var _ = Describe("Error", func() {
	It("creates an error successfully", func() {
		err := flaw.Errorf("oh no")
		Expect(err).To(MatchError("message: oh no"))
	})

	It("wraps an error successfully", func() {
		err := flaw.Wrap(fmt.Errorf("oh no"))
		Expect(err).To(MatchError("cause: oh no"))
		Expect(err.Unwrap()).To(MatchError("oh no"))
	})

	Describe("WithCode", func() {
		It("creates an error successfully", func() {
			err := flaw.Errorf("oh no").WithCode(200)
			Expect(err.Error()).To(HavePrefix("code: 200 message: oh no"))
			Expect(err.Code()).To(Equal(200))
		})

		It("returns the code", func() {
			err := flaw.Errorf("oh no").WithCode(200)
			Expect(flaw.Code(err)).To(Equal(200))
		})

		Context("when the error does not have code", func() {
			It("returns no code", func() {
				Expect(flaw.Code(fmt.Errorf("oh no"))).To(Equal(0))
			})
		})
	})

	Describe("WithStatus", func() {
		It("creates an error successfully", func() {
			err := flaw.Errorf("oh no").WithStatus(200)
			Expect(err.Status()).To(Equal(200))
		})

		It("returns the status", func() {
			err := flaw.Errorf("oh no").WithStatus(200)
			Expect(flaw.Status(err)).To(Equal(200))
		})

		Context("when the error does not have status", func() {
			It("returns no status", func() {
				Expect(flaw.Status(fmt.Errorf("oh no"))).To(Equal(0))
			})
		})
	})

	Describe("WithMessage", func() {
		It("creates an error successfully", func() {
			err := flaw.Wrap(fmt.Errorf("oh no")).WithMessage("failed")
			Expect(err.Error()).To(HavePrefix("message: failed cause: oh no"))
			Expect(err.Unwrap()).To(MatchError("oh no"))
			Expect(flaw.Message(err)).To(Equal("failed"))
		})

		Context("when the error does not have a message", func() {
			It("returns an empty message", func() {
				Expect(flaw.Message(fmt.Errorf("oh no"))).To(BeEmpty())
			})
		})
	})

	Describe("WithDetails", func() {
		It("creates an error successfully", func() {
			err := flaw.Errorf("oh no").WithDetails("some more details")
			Expect(flaw.Details(err)).To(ContainElement("some more details"))
		})

		Context("when the error does not have details", func() {
			It("returns no details", func() {
				Expect(flaw.Details(fmt.Errorf("oh no"))).To(BeEmpty())
			})
		})
	})

	Describe("WithError", func() {
		It("creates an error successfully", func() {
			err := flaw.Errorf("failed").WithError(fmt.Errorf("oh no"))
			Expect(err.Error()).To(HavePrefix("message: failed cause: oh no"))
			Expect(err.Unwrap()).To(MatchError("oh no"))
			Expect(flaw.Cause(err)).To(MatchError("oh no"))
			Expect(err.StackTrace()).NotTo(BeEmpty())
		})

		Context("when the error is not causer", func() {
			It("returns nil cause", func() {
				Expect(flaw.Cause(fmt.Errorf("oh no"))).To(MatchError("oh no"))
			})
		})
	})

	Describe("WithContext", func() {
		It("creates an error successfully", func() {
			err := flaw.Wrap(fmt.Errorf("oh no")).WithContext(flaw.Map{"user": "root"})
			Expect(flaw.Context(err)).To(HaveKeyWithValue("user", "root"))
		})

		Context("when the error does not have context", func() {
			It("returns nil context", func() {
				Expect(flaw.Context(fmt.Errorf("oh no"))).To(BeEmpty())
			})
		})

		Context("when the context is nil", func() {
			It("creates an error successfully", func() {
				err := flaw.Wrap(fmt.Errorf("oh no")).WithContext(nil)
				Expect(err.Context()).To(HaveKeyWithValue("error_cause", "oh no"))
			})
		})
	})

	Describe("Format", func() {
		It("prints the error successfully", func() {
			err := flaw.Errorf("failed").WithCode(404).WithError(fmt.Errorf("oh no"))
			Expect(fmt.Sprintf("%v", err)).To(Equal("code: 404 message: failed cause: oh no"))
		})

		Context("when the verbose printing is used", func() {
			It("prints the error successfully", func() {
				err := flaw.Errorf("failed").WithCode(404).WithError(fmt.Errorf("oh no"))
				Expect(fmt.Sprintf("%+v", err)).To(HavePrefix("    code: 404\n message: failed\n   cause: oh no\n   stack: \n"))
			})
		})
	})

	Describe("MarshalJSON", func() {
		It("marshals the error successfully", func() {
			errx := flaw.Errorf("oh no").WithCode(200)
			errx.Wrap(fmt.Errorf("failed"))

			data, err := json.Marshal(errx)
			Expect(err).To(BeNil())

			Expect(string(data)).To(Equal(`{"error_cause":"failed","error_code":200,"error_message":"oh no"}`))
		})

		Context("when the wrapped error implements MarshalJSON", func() {
			It("marshals the error successfully", func() {
				errx := flaw.Errorf("oh no").WithCode(200)
				errx.Wrap(flaw.Errorf("failed"))

				data, err := json.Marshal(errx)
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`{"error_cause":{"error_message":"failed"},"error_code":200,"error_message":"oh no"}`))
			})
		})
	})
})

var _ = Describe("ErrorCollection", func() {
	It("creates an error successfully", func() {
		errs := flaw.ErrorCollector{}
		errs = append(errs, fmt.Errorf("oh no"))
		Expect(errs).To(MatchError("[oh no]"))
	})

	Context("when the collector has more than one error", func() {
		It("creates an error successfully", func() {
			errs := flaw.ErrorCollector{}
			errs = append(errs, fmt.Errorf("oh no"))
			errs = append(errs, fmt.Errorf("oh yes"))
			Expect(errs).To(MatchError("[oh no, oh yes]"))
		})
	})

	Describe("Is", func() {
		It("returns true", func() {
			err := fmt.Errorf("not found")

			errs := flaw.ErrorCollector{}
			errs = append(errs, err)

			Expect(errors.Is(errs, err)).To(BeTrue())
		})

		Context("when the target is collector", func() {
			It("returns true", func() {
				child := fmt.Errorf("oh no")

				target := flaw.ErrorCollector{}
				target = append(target, child)

				errs := flaw.ErrorCollector{}
				errs = append(errs, child)

				Expect(errors.Is(errs, target)).To(BeTrue())
			})

			Context("when the target has more items", func() {
				It("returns false", func() {
					child := fmt.Errorf("oh no")

					target := flaw.ErrorCollector{}
					target = append(target, child)
					target = append(target, fmt.Errorf("oh yes"))

					errs := flaw.ErrorCollector{}
					errs = append(errs, child)

					Expect(errors.Is(errs, target)).To(BeFalse())
				})
			})
		})

		Context("when the error cannot be found", func() {
			It("returns false", func() {
				err := fmt.Errorf("not found")

				errs := flaw.ErrorCollector{}
				errs = append(errs, fmt.Errorf("oh no"))

				Expect(errors.Is(errs, err)).To(BeFalse())
			})
		})
	})

	Describe("As", func() {
		It("returns true", func() {
			var err *flaw.Error

			errs := flaw.ErrorCollector{}
			errs = append(errs, fmt.Errorf("oh no"))
			errs = append(errs, flaw.Errorf("not found"))

			Expect(errors.As(errs, &err)).To(BeTrue())
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(HavePrefix("message: not found"))
		})

		Context("when the error cannot be found", func() {
			It("returns false", func() {
				var err *flaw.Error

				errs := flaw.ErrorCollector{}
				errs = append(errs, fmt.Errorf("oh no"))
				errs = append(errs, fmt.Errorf("oh yes"))

				Expect(errors.As(errs, &err)).To(BeFalse())
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Wrap", func() {
		It("wraps the errors", func() {
			errs := flaw.ErrorCollector{}
			errs.Wrap(fmt.Errorf("oh no"))
			errs.Wrap(fmt.Errorf("oh yes"))

			Expect(errs).To(HaveLen(2))
			Expect(errs).To(ContainElement(fmt.Errorf("oh no")))
			Expect(errs).To(ContainElement(fmt.Errorf("oh yes")))
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
				It("unwraps the errors as nil", func() {
					errs := flaw.ErrorCollector{}
					errs = append(errs, fmt.Errorf("oh no"))
					errs = append(errs, fmt.Errorf("oh yes"))
					Expect(errs.Unwrap()).To(BeNil())
				})
			})
		})
	})

	Describe("Format", func() {
		It("prints the error successfully", func() {
			errs := flaw.ErrorCollector{}
			errs = append(errs, fmt.Errorf("oh no"))
			errs = append(errs, fmt.Errorf("oh yes"))
			Expect(fmt.Sprintf("%v", errs)).To(Equal("[oh no, oh yes]"))
		})

		Context("when the %s format is used", func() {
			It("prints the error successfully", func() {
				errs := flaw.ErrorCollector{}
				errs = append(errs, fmt.Errorf("oh no"))
				errs = append(errs, fmt.Errorf("oh yes"))
				Expect(fmt.Sprintf("%s", errs)).To(Equal("[oh no, oh yes]"))
			})
		})

		Context("when the %#v format is used", func() {
			It("prints the error successfully", func() {
				errs := flaw.ErrorCollector{}
				errs = append(errs, fmt.Errorf("oh no"))
				errs = append(errs, fmt.Errorf("oh yes"))
				Expect(fmt.Sprintf("%#v", errs)).To(Equal("[]error{oh no, oh yes}"))
			})
		})

		Context("when the %v format is used", func() {
			It("prints the error successfully", func() {
				errs := flaw.ErrorCollector{}
				errs = append(errs, fmt.Errorf("oh no"))
				errs = append(errs, fmt.Errorf("oh yes"))
				Expect(fmt.Sprintf("%+v", errs)).To(Equal(" --- oh no\n --- oh yes"))
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
				errs = append(errs, flaw.Errorf("oh no"))

				data, err := json.Marshal(errs)
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`[{"error_message":"oh no"}]`))
			})
		})
	})
})

var _ = Describe("ErrorConstant", func() {
	It("creates a error constant successfully", func() {
		const err = flaw.ErrorConstant("EOF")
		Expect(err).To(MatchError("EOF"))
	})
})
