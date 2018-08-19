/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package master

import (
	samplev1alpha1 "github.com/seanchann/sample-apiserver/pkg/apis/sample/v1alpha1"
	samplerest "github.com/seanchann/sample-apiserver/pkg/registry/sample/rest"

	"github.com/seanchann/apimaster/pkg/apiserver"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

//DefaultRESTStorageProvider build default install api rest provider
func DefaultRESTStorageProvider() []apiserver.RESTStorageProvider {
	return []apiserver.RESTStorageProvider{
		samplerest.RESTStorageProvider{},
	}
}

//DefaultAPIResourceConfigSource default api resource config
func DefaultAPIResourceConfigSource() *serverstorage.ResourceConfig {
	ret := serverstorage.NewResourceConfig()
	// NOTE: GroupVersions listed here will be enabled by default. Don't put alpha versions in the list.
	ret.EnableVersions(
		samplev1alpha1.SchemeGroupVersion,
	)
	// disable alpha versions explicitly so we have a full list of what's possible to serve
	ret.DisableVersions()

	return ret
}
