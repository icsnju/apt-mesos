package impl

import "github.com/icsnju/apt-mesos/registry"

func (core *Core) getAgentPort() string {
	return "18080"
}

func (core *Core) checkAgentHealth(node *registry.Node) bool {
	return true
}
