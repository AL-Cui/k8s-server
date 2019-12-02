package modules

import (
	"k8s-server/modules/pod"
	"k8s-server/utils/errors"
	"k8s-server/models"
)
var KubernetesServer *Backend

type Backend struct {
	DB models.Model
	PodManager *pod.Manager
	inited    bool
}

func NewBackend() (*Backend, error) {
	if KubernetesServer != nil && KubernetesServer.inited {
		return KubernetesServer, nil
	}
	podManager,err := pod.NewManager()
	if err != nil {
		return nil, errors.Wrap(err, def.ErrPodModule,"init pod module failed")
	}
	m, err := models.GetModel()
	if err != nil {
		return nil, errors.Wrap(err, def.ErrGeneralDBConnect,
			"init database failed")
	}
	backend := &Backend{
		DB:             m,
		PodManager: podManager,
		inited: true,
	}
}