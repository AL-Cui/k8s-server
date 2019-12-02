package controllers

import (
	"strings"
	"k8s-server/modulers"
)

type Pod Struct {
	BaseController
}

func (p *Pod) nestPrepare() {

}

func (p *Pod) ListPodsFromNamespace() {
	pods, err := modulers.KubernetesServer.PodManager
}