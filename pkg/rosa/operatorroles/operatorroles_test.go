package operatorroles_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-online/ocm-common/pkg/rosa/operatorroles"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var _ = Describe("Operator Role functions", func() {
	var _ = Describe("Validates GetOperatorRolesArnsMap function", func() {
		It("Checks the operator roles are retrieved from a cluster", func() {
			fakeCluster, err := cmv1.NewCluster().
				AWS(
					cmv1.NewAWS().
						STS(
							cmv1.NewSTS().
								OperatorIAMRoles(
									cmv1.NewOperatorIAMRole().Name("operator-role-1").RoleARN("arn-1"),
									cmv1.NewOperatorIAMRole().Name("operator-role-2").RoleARN("arn-2"),
									cmv1.NewOperatorIAMRole().Name("operator-role-3").RoleARN("arn-3"),
								),
						),
				).Build()
			Expect(err).ToNot(HaveOccurred())
			operatorRolesMap := GetOperatorRolesArnsMap(fakeCluster)
			Expect(len(operatorRolesMap)).To(Equal(3))
			for name, arn := range operatorRolesMap {
				Expect(name).ToNot(Equal(""))
				Expect(arn).ToNot(Equal(""))
			}
		})
	})
})
