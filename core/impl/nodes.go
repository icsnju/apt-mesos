package impl

import (
	"errors"
	"time"

	"github.com/icsnju/apt-mesos/mesosproto"
	"github.com/icsnju/apt-mesos/registry"
	scheduler "github.com/icsnju/apt-mesos/scheduler/impl"
)

// ErrNodeNotExists defined errors
var (
	ErrNodeNotExists = errors.New("Specific task not exist")
)

// RegisterNode register node to registry
func (core *Core) RegisterNode(id string, node *registry.Node) error {
	return core.nodes.Add(id, node)
}

// GetNode return node with id
func (core *Core) GetNode(id string) (*registry.Node, error) {
	if node := core.nodes.Get(id); node != nil {
		return node.(*registry.Node), nil
	}
	return nil, ErrNodeNotExists
}

// UpdateNode update node information
func (core *Core) UpdateNode(id string, node *registry.Node) error {
	node.LastUpdateTime = time.Now().Unix()
	if err := core.nodes.Update(id, node); err != nil {
		return err
	}
	return nil
}

// DeleteNode delete node with given id
func (core *Core) DeleteNode(id string) error {
	return core.nodes.Delete(id)
}

// ExistsNode return if node exists in registry
func (core *Core) ExistsNode(id string) bool {
	return core.nodes.Exists(id)
}

// GetAllNodes return all nodes
func (core *Core) GetAllNodes() []*registry.Node {
	rawList := core.nodes.List()
	nodes := make([]*registry.Node, len(rawList))

	for i, v := range rawList {
		nodes[i] = v.(*registry.Node)
	}
	return nodes
}

func (core *Core) updateNodesByOffer(offers []*mesosproto.Offer) {
	for _, offer := range offers {
		// update resources
		resources := make(map[string]*mesosproto.Resource)
		for _, resource := range offer.GetResources() {
			resources[resource.GetName()] = resource
		}

		// if it is a new node
		slaveID := offer.GetSlaveId().GetValue()
		if exists := core.ExistsNode(slaveID); !exists {
			node := &registry.Node{
				ID:             slaveID,
				Hostname:       offer.GetHostname(),
				LastUpdateTime: time.Now().Unix(),
				Resources:      resources,
			}
			core.RegisterNode(slaveID, node)
		} else {
			node, _ := core.GetNode(slaveID)
			node.Resources = resources
			node.LastUpdateTime = time.Now().Unix()
			core.UpdateNode(slaveID, node)
		}
	}
}

func (core *Core) updateNodeByTask(id string, task *registry.Task) {
	node, _ := core.GetNode(id)
	for _, resource := range task.Resources {
		// Update scalar
		if resource.GetType().String() == "SCALAR" {
			newScalar := node.Resources[resource.GetName()].GetScalar().GetValue() - resource.GetScalar().GetValue()
			node.Resources[resource.GetName()].Scalar.Value = &newScalar
		} else if resource.GetType().String() == "RANGES" {
			// Update ranges
			node.Resources[resource.GetName()].Ranges = scheduler.RangeUsedUpdate(resource.GetRanges(), node.Resources[resource.GetName()].GetRanges())
		}
	}
	core.UpdateNode(id, node)
}
