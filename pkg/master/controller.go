package master

import (
	"github.com/golang/glog"
	"github.com/seanchann/apimaster/pkg/apiserver"
	"github.com/seanchann/apimaster/pkg/util/async"

	genericapiserver "k8s.io/apiserver/pkg/server"
)

//Controller controller
type Controller struct {
	runner *async.Runner
}

//NewController create controller
func NewController(para []interface{}) (*apiserver.ControllerProvider, error) {
	ctrl := &Controller{}
	for _, v := range para {
		glog.Infof("input parameter %+v", v)
	}

	return &apiserver.ControllerProvider{
		NameFunc:        ctrl.Name,
		PostFunc:        ctrl.PostStartHook,
		PreShutdownFunc: ctrl.PreShutdownHook,
	}, nil
}

//Name return controller name
func (c *Controller) Name() string {
	return "sample-controller"
}

//PostStartHook when apiserver startup call this hook
func (c *Controller) PostStartHook(hookContext genericapiserver.PostStartHookContext) error {
	c.Start()
	return nil
}

//PreShutdownHook when apiserver shutdown call this hook
func (c *Controller) PreShutdownHook() error {
	c.Stop()
	return nil
}

// Start begins the core controller loops that must exist for bootstrapping
// a cluster.
func (c *Controller) Start() {
	if c.runner != nil {
		return
	}
}

//Stop stop all runner
func (c *Controller) Stop() {
	glog.Info("shutdown http service")
}
