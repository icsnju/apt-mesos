package registry

import "github.com/icsnju/apt-mesos/mesosproto"

// Node is one node of cluster
type Node struct {
	ID             string                          `json:"id"`
	Hostname       string                          `json:"hostname"`
	NodeType       string                          `json:"node_type"`
	Resources      map[string]*mesosproto.Resource `json:"resources"`
	LastUpdateTime int64                           `json:"last_update_time"`
}
