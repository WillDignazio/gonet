package gonet

import (
	"errors"
	"sync"
)

type GoNetCtx struct {
	rwlock *sync.RWMutex
	interfaces map[string]GoNetInterface
	config *GoNetConfig
}

func (ctx *GoNetCtx) Interfaces() []GoNetInterface {
	ctx.rwlock.RLock()
	defer ctx.rwlock.RUnlock()

	amnt := len(ctx.interfaces)
	ret := make([]GoNetInterface, amnt)

	idx := 0
	for _, val := range ctx.interfaces {
		ret[idx] = val
		idx += 1
	}

	return ret
}

func (ctx *GoNetCtx) GetInterface(name string) (GoNetInterface, error) {
	ctx.rwlock.RLock()
	defer ctx.rwlock.RUnlock()
	
	iface, present := ctx.interfaces[name]
	if !present {
		return nil, errors.New("Interface not found")
	}

	return iface, nil
}

func (ctx *GoNetCtx) AddInterface(iface GoNetInterface) error {
	if iface == nil {
		return errors.New("Attempted to add NIL interface")
	}

	ctx.rwlock.Lock()
	defer ctx.rwlock.Unlock()

	name := iface.Name()
	if _, present := ctx.interfaces[name]; present {
		return errors.New("Interface with same name already present")
	}

	ctx.interfaces[name] = iface
	return nil
}

func (ctx *GoNetCtx) RemoveInterface(name string) error {
	ctx.rwlock.Lock()
	defer ctx.rwlock.Unlock()

	if _, present := ctx.interfaces[name]; !present {
		return errors.New("Attempted to remove interface outside range of available interfaces")
	}

	delete(ctx.interfaces, name)
	return nil
}

func (ctx *GoNetCtx) Config() *GoNetConfig {
	return ctx.config
}

func NewGoNetCtx(config *GoNetConfig) (*GoNetCtx) {
	return &GoNetCtx{
		rwlock: new(sync.RWMutex),
		interfaces: make(map[string]GoNetInterface, 0),
		config: config,
	}
}
