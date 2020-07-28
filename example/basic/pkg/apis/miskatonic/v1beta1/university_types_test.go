/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	. "sigs.k8s.io/apiserver-builder-alpha/example/basic/pkg/apis/miskatonic/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("University", func() {
	var instance University
	var expected University

	BeforeEach(func() {
		instance = University{}
		instance.Name = "miskatonic-university"
		instance.Namespace = "default"
		instance.Spec.FacultySize = 7
		instance.Spec.ServiceSpec = corev1.ServiceSpec{}
		instance.Spec.ServiceSpec.ClusterIP = "1.1.1.1"

		expected = instance
		val := 15
		expected.Spec.MaxStudents = &val
		expected.Spec.ServiceSpec = corev1.ServiceSpec{}
		expected.Spec.ServiceSpec.ClusterIP = "1.1.1.1"
	})

	AfterEach(func() {
		cs.Delete(context.TODO(), &instance)
	})

	Describe("when sending a storage request", func() {
		Context("for a valid config", func() {
			It("should provide CRUD access to the object", func() {

				By("returning success from the create request")
				actual := instance.DeepCopy()
				err := cs.Create(context.TODO(), actual)
				Expect(err).ShouldNot(HaveOccurred())

				By("defaulting the expected fields")
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("returning the item for list requests")

				var result UniversityList
				err = cs.List(context.TODO(), &result)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected.Spec))

				By("returning the item for get requests")
				err = cs.Get(context.TODO(), client.ObjectKey{
					Namespace: instance.Namespace,
					Name:      instance.Name,
				}, actual)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("deleting the item for delete requests")
				err = cs.Delete(context.TODO(), &instance)
				Expect(err).ShouldNot(HaveOccurred())
				err = cs.List(context.TODO(), &result)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(0))
			})
		})
		Context("for an invalid config", func() {
			It("should fail if there are too many students", func() {
				val := 151
				instance.Namespace = "university-test-too-many"
				instance.Spec.MaxStudents = &val
				err := cs.Create(context.TODO(), &instance)
				Expect(err).Should(HaveOccurred())
			})

			It("should fail if there are not enough students", func() {
				val := 0
				instance.Namespace = "university-test-not-enough"
				instance.Spec.MaxStudents = &val
				err := cs.Create(context.TODO(), &instance)
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
