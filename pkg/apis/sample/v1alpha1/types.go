/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TestSpec struct {
	//Family is the harbor of love
	Family string `json:"family"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//Test is sample struct to use our apimaster
type Test struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TestSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//TestList is sample list struct to use our apimaster
type TestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Test
}

//UserInfo is a mysql users map
type UserInfo struct {
}

type UserSpec struct {
	ID          int64  `json:"-" freezer:"column:id;const" gorm:"column:id"`
	Passwd      string `json:"passwd,omitempty" freezer:"column:passwd" gorm:"column:passwd"`
	Email       string `json:"email,omitempty" freezer:"column:email" gorm:"column:email"`
	Name        string `json:"name" freezer:"column:user_name;resoucekey" gorm:"column:user_name"`
	EmailVerify bool   `json:"emailVerify,omitempty" freezer:"column:is_email_verify" gorm:"column:is_email_verify"`
	Status      bool   `json:"status,omitempty" freezer:"column:status" gorm:"column:status"`
	RawObj      []byte `json:"-" freezer:"column:rawobj" gorm:"column:rawobj"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec,omitempty"  freezer:"table:user"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []User `json:"spec,omitempty"`
}
