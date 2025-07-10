package kernel

import (
	"fmt"
	"sync"
)

type Factory func(k IKernel) (interface{}, error)

type IKernel interface {
	Register(key string, factory Factory)
	RegisterSingleton(key string, factory Factory)
	Get(key string) (interface{}, error)
}

type kernel struct {
	mu         sync.RWMutex
	bindings   map[string]Factory
	singletons map[string]interface{}
}

func New() IKernel {
	return &kernel{
		bindings:   make(map[string]Factory),
		singletons: make(map[string]interface{}),
	}
}

func (k *kernel) Register(key string, factory Factory) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.bindings[key] = factory
}

func (k *kernel) RegisterSingleton(key string, factory Factory) {
	k.Register(key, func(kern IKernel) (interface{}, error) {
		k := kern.(*kernel)
		k.mu.RLock()
		instance, ok := k.singletons[key]
		k.mu.RUnlock()
		if ok {
			return instance, nil
		}
		k.mu.Lock()
		defer k.mu.Unlock()
		if instance, ok := k.singletons[key]; ok {
			return instance, nil
		}
		newInstance, err := factory(k)
		if err != nil {
			return nil, err
		}
		k.singletons[key] = newInstance
		return newInstance, nil
	})
}

func (k *kernel) Get(key string) (interface{}, error) {
	k.mu.RLock()
	factory, ok := k.bindings[key]
	k.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("kernel: service with key '%s' not found", key)
	}
	return factory(k)
}