/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package app

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/seanchann/sample-apiserver/pkg/app/options"
	clientset "github.com/seanchann/sample-apiserver/pkg/client/clientset/internalversion"
	informers "github.com/seanchann/sample-apiserver/pkg/client/informers/internalversion"
	"github.com/seanchann/sample-apiserver/pkg/master"
	utilflag "github.com/seanchann/sample-apiserver/pkg/util/flag"

	"github.com/seanchann/apimaster/pkg/api/legacyscheme"
	"github.com/seanchann/apimaster/pkg/apiserver"
	insecureserver "github.com/seanchann/apimaster/pkg/apiserver/server"

	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	genericapiserver "k8s.io/apiserver/pkg/server"
	apiserverstorage "k8s.io/apiserver/pkg/server/storage"
	clientgoinformers "k8s.io/client-go/informers"
	clientgoclientset "k8s.io/client-go/kubernetes"
)

const defaultEtcdPathPrefix = "/registry/ast.cticlound.cn"

func init() {
	fmt.Printf("in app init\r\n")
}

// NewASTAgentCommand provides a CLI handler for 'start master' command
// with a default ASTAgentOptions.
func NewASTAgentCommand(stopCh <-chan struct{}) *cobra.Command {
	s := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Short: "Launch a astagent server",
		Long:  "Launch a astagent server",
		RunE: func(cmd *cobra.Command, args []string) error {
			utilflag.PrintFlags(cmd.Flags())

			// set default options
			completedOptions, err := Complete(s)
			if err != nil {
				return err
			}

			// validate options
			if errs := completedOptions.Validate(); len(errs) != 0 {
				return utilerrors.NewAggregate(errs)
			}

			return Run(completedOptions, stopCh)
		},
	}
	s.AddFlags(cmd.Flags())

	return cmd
}

// Run runs the specified APIServer.  This should never exit.
func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	// glog.Infof("Version: %+v", version.Get())

	server, err := CreateServerChain(completeOptions, stopCh)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run(stopCh)
}

// CreateAPIServer creates and wires a workable agent-apiserver
func CreateAPIServer(apiServerConfig *apiserver.Config,
	delegateAPIServer genericapiserver.DelegationTarget,
	sharedInformers informers.SharedInformerFactory,
	versionedInformers clientgoinformers.SharedInformerFactory) (*apiserver.APIServer, error) {
	apiServer, err := apiServerConfig.Complete(versionedInformers).New(delegateAPIServer)
	if err != nil {
		return nil, err
	}

	return apiServer, nil
}

// CreateServerChain creates the apiservers connected via delegation.
func CreateServerChain(completedOptions completedServerRunOptions, stopCh <-chan struct{}) (*genericapiserver.GenericAPIServer, error) {

	apiServerCfg, sharedInformers, versionedInformers, insecureServingOptions, err := CreateAPIServerConfig(completedOptions)
	if err != nil {
		return nil, err
	}

	apiServer, err := CreateAPIServer(apiServerCfg, genericapiserver.NewEmptyDelegate(), sharedInformers, versionedInformers)
	if err != nil {
		return nil, err
	}

	if insecureServingOptions != nil {
		insecureHandlerChain := insecureserver.BuildInsecureHandlerChain(apiServer.GenericAPIServer.UnprotectedHandler(), apiServerCfg.GenericConfig)
		if err := insecureserver.NonBlockingRun(insecureServingOptions, insecureHandlerChain, apiServerCfg.GenericConfig.RequestTimeout, stopCh); err != nil {
			return nil, err
		}
	}

	return apiServer.GenericAPIServer, nil
}

// CreateAPIServerConfig creates all the resources for running the API server, but runs none of them
func CreateAPIServerConfig(s completedServerRunOptions) (
	config *apiserver.Config,
	sharedInformers informers.SharedInformerFactory,
	versionedInformers clientgoinformers.SharedInformerFactory,
	insecureServingInfo *insecureserver.InsecureServingInfo,
	lastErr error,
) {

	genericConfig := genericapiserver.NewConfig(legacyscheme.Codecs)
	if lastErr = s.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if insecureServingInfo, lastErr = s.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = s.SecureServing.ApplyTo(&genericConfig.SecureServing, &genericConfig.LoopbackClientConfig); lastErr != nil {
		return
	}
	if lastErr = s.Authentication.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = s.Audit.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = s.Features.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = s.APIEnablement.ApplyTo(genericConfig, master.DefaultAPIResourceConfigSource(), legacyscheme.Scheme); lastErr != nil {
		return
	}

	genericConfig.SwaggerConfig = genericapiserver.DefaultSwaggerConfig()
	genericConfig.SwaggerConfig.SwaggerFilePath = s.SwaggerUIFilePath

	storageFactory, lastErr := BuildStorageFactory(s.ServerRunOptions, genericConfig.MergedResourceConfig)
	if lastErr != nil {
		return
	}
	if lastErr = s.Mysql.ApplyWithStorageFactoryTo(storageFactory, genericConfig); lastErr != nil {
		return
	}

	// Use protobufs for self-communication.
	// Since not every generic apiserver has to support protobufs, we
	// cannot default to it in generic apiserver and need to explicitly
	// set it in kube-apiserver.
	genericConfig.LoopbackClientConfig.ContentConfig.ContentType = "application/vnd.kubernetes.protobuf"
	genericConfig.LoopbackClientConfig.Timeout = 10 * time.Minute

	client, err := clientset.NewForConfig(genericConfig.LoopbackClientConfig)
	if err != nil {
		lastErr = fmt.Errorf("failed to create clientset: %v", err)
		return
	}

	clientConfig := genericConfig.LoopbackClientConfig
	sharedInformers = informers.NewSharedInformerFactory(client, genericConfig.LoopbackClientConfig.Timeout)
	clientgoExternalClient, err := clientgoclientset.NewForConfig(clientConfig)
	if err != nil {
		lastErr = fmt.Errorf("failed to create real external clientset: %v", err)
		return
	}
	versionedInformers = clientgoinformers.NewSharedInformerFactory(clientgoExternalClient, genericConfig.LoopbackClientConfig.Timeout)

	config = &apiserver.Config{
		GenericConfig: genericConfig,
		ExtraConfig: apiserver.ExtraConfig{
			APIResourceConfigSource: master.DefaultAPIResourceConfigSource(),
			RESTStorageProviders:    master.DefaultRESTStorageProvider(),
		},
	}

	return
}

// completedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type completedServerRunOptions struct {
	*options.ServerRunOptions
}

// Complete set default ServerRunOptions.
// Should be called after server flags parsed.
func Complete(s *options.ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions
	options.ServerRunOptions = s
	return options, nil
}

// BuildStorageFactory constructs the storage factory. If encryption at rest is used, it expects
// all supported KMS plugins to be registered in the KMS plugin registry before being called.
func BuildStorageFactory(s *options.ServerRunOptions, apiResourceConfig *apiserverstorage.ResourceConfig) (*apiserverstorage.DefaultStorageFactory, error) {
	storageGroupsToEncodingVersion, err := s.StorageSerialization.StorageGroupsToEncodingVersion()
	if err != nil {
		return nil, fmt.Errorf("error generating storage version map: %s", err)
	}
	storageFactory, err := apiserver.NewStorageFactory(
		s.Mysql.StorageConfig, s.Mysql.DefaultStorageMediaType, legacyscheme.Codecs,
		apiserverstorage.NewDefaultResourceEncodingConfig(legacyscheme.Scheme), storageGroupsToEncodingVersion,
		// The list includes resources that need to be stored in a different
		// group version than other resources in the groups.
		// FIXME (soltysh): this GroupVersionResource override should be configurable
		[]schema.GroupVersionResource{},
		apiResourceConfig)
	if err != nil {
		return nil, fmt.Errorf("error in initializing storage factory: %s", err)
	}

	return storageFactory, nil
}
