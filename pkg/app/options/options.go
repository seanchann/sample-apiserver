/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package options

import (
	"github.com/spf13/pflag"

	apioptions "github.com/seanchann/apimaster/pkg/apiserver/options"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/storage/storagebackend"
)

//ServerRunOptions astagent custom options
type ServerRunOptions struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions
	Etcd                    *genericoptions.EtcdOptions
	SecureServing           *genericoptions.SecureServingOptionsWithLoopback
	InsecureServing         *apioptions.InsecureServingOptions
	Audit                   *genericoptions.AuditOptions
	Features                *genericoptions.FeatureOptions
	Admission               *apioptions.AdmissionOptions
	Authentication          *apioptions.BuiltInAuthenticationOptions
	Authorization           *apioptions.BuiltInAuthorizationOptions
	StorageSerialization    *apioptions.StorageSerializationOptions
	APIEnablement           *genericoptions.APIEnablementOptions

	//our custom options
	SwaggerUIFilePath string
}

//NewServerRunOptions new a ASTAgentOptions
func NewServerRunOptions() *ServerRunOptions {
	o := &ServerRunOptions{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		Etcd:                 genericoptions.NewEtcdOptions(storagebackend.NewDefaultConfig(apioptions.DefaultEtcdPathPrefix, nil)),
		SecureServing:        apioptions.NewSecureServingOptions(),
		InsecureServing:      apioptions.NewInsecureServingOptions(),
		Audit:                genericoptions.NewAuditOptions(),
		Features:             genericoptions.NewFeatureOptions(),
		Admission:            apioptions.NewAdmissionOptions(),
		Authentication:       apioptions.NewBuiltInAuthenticationOptions().WithAll(),
		Authorization:        apioptions.NewBuiltInAuthorizationOptions(),
		StorageSerialization: apioptions.NewStorageSerializationOptions(),
		APIEnablement:        genericoptions.NewAPIEnablementOptions(),
	}

	return o
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet
func (o *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	o.GenericServerRunOptions.AddUniversalFlags(fs)
	o.Etcd.AddFlags(fs)
	o.SecureServing.AddFlags(fs)
	o.InsecureServing.AddFlags(fs)
	o.InsecureServing.AddDeprecatedFlags(fs)
	o.Audit.AddFlags(fs)
	o.Features.AddFlags(fs)
	o.Authentication.AddFlags(fs)
	o.Authorization.AddFlags(fs)
	o.StorageSerialization.AddFlags(fs)
	o.APIEnablement.AddFlags(fs)
	o.Admission.AddFlags(fs)

	fs.StringVar(&o.SwaggerUIFilePath, "swagger-ui-file-path", o.SwaggerUIFilePath, ""+
		"the swagger-ui file path .default /swagger-ui/")
}
