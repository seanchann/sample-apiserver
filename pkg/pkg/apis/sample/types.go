/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package sample

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type TestSpec struct {
	Family string
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Test struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec TestSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TestList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Test
}
