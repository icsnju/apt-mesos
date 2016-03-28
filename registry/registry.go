package registry

import (
	"errors"
	"sync"
)

// Registry to manager tasks
type Registry struct {
	sync.RWMutex
	items map[string]interface{}
}

// NewRegistry instantiate and return a new Registry
func NewRegistry() *Registry {
	return &Registry{
		items: make(map[string]interface{}),
	}
}

// Add is called when user submit a task and add the task to the registry
func (registry *Registry) Add(id string, item interface{}) error {
	registry.Lock()
	defer registry.Unlock()

	registry.items[id] = item
	return nil
}

// Exists return if registry has the id
func (registry *Registry) Exists(id string) bool {
	_, exists := registry.items[id]
	return exists
}

// Get : Get the task that specified id
func (registry *Registry) Get(id string) interface{} {
	registry.RLock()
	defer registry.RUnlock()

	item, exists := registry.items[id]
	if !exists {
		return nil
	}

	return item
}

// List Return all the tasks in registry
func (registry *Registry) List() []interface{} {
	registry.RLock()
	defer registry.RUnlock()

	var i = 0
	result := make([]interface{}, len(registry.items))

	for _, v := range registry.items {
		result[i] = v
		i++
	}

	return result
}

// Delete Give an id and delete the task
func (registry *Registry) Delete(id string) error {
	registry.Lock()
	defer registry.Unlock()

	_, exists := registry.items[id]
	if !exists {
		return errors.New("Registry delete error")
	}

	delete(registry.items, id)
	return nil
}

// Update update a task which give the specfic string and a new struct
func (registry *Registry) Update(id string, item interface{}) error {
	registry.Lock()
	defer registry.Unlock()

	_, exists := registry.items[id]
	if !exists {
		return errors.New("Registry update error")
	}

	registry.items[id] = item
	return nil
}
