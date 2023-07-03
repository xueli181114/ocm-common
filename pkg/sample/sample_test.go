package sample_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-common/pkg/sample"
)

var _ = Describe("Sample", func() {

	Context("Sample Test", func() {
		Context("Add Test", func() {
			It("Test positive output of two positive numbers", func() {
				Expect(sample.Add(10, 10)).To(Equal(20))
				Expect(sample.Add(200, 200)).To(Equal(400))
				Expect(sample.Add(4000, 4000)).To(Equal(8000))
			})

			It("Test positive output of two negative numbers", func() {
				Expect(sample.Add(-10, -10)).To(Equal(-20))
				Expect(sample.Add(-200, -200)).To(Equal(-400))
				Expect(sample.Add(-4000, -4000)).To(Equal(-8000))
			})
		})

		Context("Subtract Test", func() {
			It("Test positive output of two positive numbers", func() {
				Expect(sample.Subtract(10, 5)).To(Equal(5))
				Expect(sample.Subtract(200, 150)).To(Equal(50))
				Expect(sample.Subtract(4000, 3500)).To(Equal(500))
			})

			It("Test positive output of two negative numbers", func() {
				Expect(sample.Subtract(-10, -10)).To(Equal(0))
				Expect(sample.Subtract(-200, -200)).To(Equal(0))
				Expect(sample.Subtract(-4000, -4000)).To(Equal(0))
			})
		})
	})
})
