/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package rest

import (
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"

	"github.com/seanchann/apimaster/pkg/api/legacyscheme"

	"github.com/seanchann/sample-apiserver/pkg/apis/sample"
	sampleapiv1alpha1 "github.com/seanchann/sample-apiserver/pkg/apis/sample/v1alpha1"
	teststore "github.com/seanchann/sample-apiserver/pkg/registry/sample/test/storage"
	userstore "github.com/seanchann/sample-apiserver/pkg/registry/sample/user/storage"
)

//RESTStorageProvider providers information needed to build RESTStorage for core.
type RESTStorageProvider struct {
}

//NewRESTStorage create a RESTStorage provider
func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(sample.GroupName,
		legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	if apiResourceConfigSource.VersionEnabled(sampleapiv1alpha1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[sampleapiv1alpha1.SchemeGroupVersion.Version] = p.v1alpha1Storage(apiResourceConfigSource, restOptionsGetter)
	}

	return apiGroupInfo, true
}

func (p RESTStorageProvider) v1alpha1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) map[string]rest.Storage {

	// test
	testStorage := teststore.NewREST()
	userStorage := userstore.NewREST(restOptionsGetter)

	restStorageMap := map[string]rest.Storage{
		"tests": testStorage,
		"users": userStorage,
	}

	return restStorageMap
}

// GroupName return ami group name
func (p RESTStorageProvider) GroupName() string {
	return sample.GroupName
}
