package validation

import (
	"github.com/seanchann/sample-apiserver/pkg/apis/sample"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateObjectMetaUpdate validates an object's metadata when updated
func ValidateObjectMetaUpdate(newMeta, oldMeta *metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(newMeta, oldMeta, fldPath)

	return allErrs
}

// ValidateTest validates a Test.
func ValidateTest(f *sample.Test) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateTestSpec(&f.Spec, field.NewPath("spec"))...)

	return allErrs
}

// ValidateTestSpec validates a TestSpec.
func ValidateTestSpec(s *sample.TestSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(s.Family) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("family"), s.Family, "cannot be empty in Test"))
	}

	return allErrs
}

// ValidateUser validates a User.
func ValidateUser(f *sample.User) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateUserSpec(&f.Spec, field.NewPath("spec"))...)

	return allErrs
}

// ValidateUserSpec validates a UserSpec.
func ValidateUserSpec(s *sample.UserSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(s.Passwd) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("id"), s.Passwd, "cannot be empty in User"))
	}

	return allErrs
}

//ValidateUserUpdate
func ValidateUserUpdate(new, old *sample.User) field.ErrorList {
	fldPath := field.NewPath("metadata")
	allErrs := ValidateObjectMetaUpdate(&new.ObjectMeta, &old.ObjectMeta, fldPath)

	return allErrs
}
