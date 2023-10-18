package validations

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWS IAM Functions", func() {
	Describe("GetRoleName", func() {
		It("should generate a role name with the given prefix and role name", func() {
			prefix := "myPrefix"
			roleName := "myRole"
			expectedName := fmt.Sprintf("%s-%s-Role", prefix, roleName)

			name := GetRoleName(prefix, roleName)

			Expect(name).To(Equal(expectedName))
		})

		It("should truncate the generated name if it exceeds 64 characters", func() {
			prefix := "myPrefix"
			roleName := "myVeryLongRoleNameThatExceedsSixtyFourCharacters123456"
			expectedName := "myPrefix-myVeryLongRoleNameThatExceedsSixtyFourCharacters123456-"

			name := GetRoleName(prefix, roleName)

			Expect(name).To(Equal(expectedName))
		})
	})

	Describe("isManagedRole", func() {
		It("should return true if the 'ManagedPolicies' tag has the value 'true'", func() {
			roleTags := []*iam.Tag{
				{Key: aws.String(ManagedPolicies), Value: aws.String("true")},
			}

			result := IsManagedRole(roleTags)

			Expect(result).To(BeTrue())
		})

		It("should return false if the 'ManagedPolicies' tag does not have the value 'true'", func() {
			roleTags := []*iam.Tag{
				{Key: aws.String(ManagedPolicies), Value: aws.String("false")},
			}

			result := IsManagedRole(roleTags)

			Expect(result).To(BeFalse())
		})

		It("should return false if the 'ManagedPolicies' tag is not present", func() {
			roleTags := []*iam.Tag{
				{Key: aws.String("SomeOtherTag"), Value: aws.String("true")},
			}

			result := IsManagedRole(roleTags)

			Expect(result).To(BeFalse())
		})
	})

	var _ = Describe("HasCompatibleVersionTags", func() {
		var iamTags []*iam.Tag

		BeforeEach(func() {
			iamTags = []*iam.Tag{
				{Key: aws.String(OpenShiftVersion), Value: aws.String("1.2.3")},
				{Key: aws.String("SomeOtherTag"), Value: aws.String("value")},
			}
		})

		It("should return true if the version tag matches the provided version", func() {
			version := "1.2.3"

			result, err := HasCompatibleVersionTags(iamTags, version)

			Expect(result).To(BeTrue())
			Expect(err).To(BeNil())
		})

		It("should return false if the version tag does not match the provided version", func() {
			version := "2.0.0"

			result, err := HasCompatibleVersionTags(iamTags, version)

			Expect(result).To(BeFalse())
			Expect(err).To(BeNil())
		})

		It("should return false if the version tag is not present", func() {
			version := "1.2.3"
			iamTags = []*iam.Tag{
				{Key: aws.String("SomeOtherTag"), Value: aws.String("value")},
			}

			result, err := HasCompatibleVersionTags(iamTags, version)

			Expect(result).To(BeFalse())
			Expect(err).To(BeNil())
		})

		It("should return an error if the provided version is not a valid semantic version", func() {
			version := "invalid-version"

			result, err := HasCompatibleVersionTags(iamTags, version)

			Expect(result).To(BeFalse())
			Expect(err).ToNot(BeNil())
		})
	})

	var _ = Describe("IamResourceHasTag", func() {
		It("should return true if the tag with the specified key and value exists", func() {
			iamTags := []*iam.Tag{
				{Key: aws.String("Tag1"), Value: aws.String("Value1")},
				{Key: aws.String("Tag2"), Value: aws.String("Value2")},
			}
			tagKey := "Tag1"
			tagValue := "Value1"

			result := IamResourceHasTag(iamTags, tagKey, tagValue)

			Expect(result).To(BeTrue())
		})

		It("should return false if the tag with the specified key and value does not exist", func() {
			iamTags := []*iam.Tag{
				{Key: aws.String("Tag1"), Value: aws.String("Value1")},
				{Key: aws.String("Tag2"), Value: aws.String("Value2")},
			}
			tagKey := "Tag3"
			tagValue := "Value3"

			result := IamResourceHasTag(iamTags, tagKey, tagValue)

			Expect(result).To(BeFalse())
		})
	})
})
