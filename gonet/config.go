package gonet

import (
	"sort"
	"sync"
)

const (
	DEBUG_KEY = "gonet.debug"
)

type GoNetConfig struct {
	rwlock  *sync.RWMutex
	ints    map[string]int
	floats  map[string]float32
	strings map[string]string
	bools   map[string]bool
}

func (conf *GoNetConfig) SetBool(key string, val bool) {
	conf.rwlock.Lock()
	defer conf.rwlock.Unlock()

	conf.bools[key] = val
}

func (conf *GoNetConfig) GetBool(key string, def bool) bool {
	conf.rwlock.RLock()

	if val, present := conf.bools[key]; present {
		conf.rwlock.RUnlock()
		return val
	}

	conf.rwlock.RUnlock()
	conf.SetBool(key, def)

	return conf.GetBool(key, def)
}

func (conf *GoNetConfig) SetInt(key string, val int) {
	conf.rwlock.Lock()
	defer conf.rwlock.Unlock()

	conf.ints[key] = val
}

func (conf *GoNetConfig) GetInt(key string, def int) int {
	conf.rwlock.RLock()

	if val, present := conf.ints[key]; present {
		conf.rwlock.RUnlock()
		return val
	}

	conf.rwlock.RUnlock()
	conf.SetInt(key, def)
	return conf.GetInt(key, def)
}

func (conf *GoNetConfig) SetFloat(key string, val float32) {
	conf.rwlock.Lock()
	defer conf.rwlock.Unlock()

	conf.floats[key] = val
}

func (conf *GoNetConfig) GetFloat(key string, def float32) float32 {
	conf.rwlock.RLock()

	if val, present := conf.floats[key]; present {
		conf.rwlock.RUnlock()
		return val
	}
	conf.rwlock.RUnlock()
	conf.SetFloat(key, def)
	return conf.GetFloat(key, def)
}

func (conf *GoNetConfig) SetString(key string, val string) {
	conf.rwlock.Lock()
	defer conf.rwlock.Unlock()

	conf.strings[key] = val
}

func (conf *GoNetConfig) GetString(key string, def string) string {
	conf.rwlock.RLock()

	if val, present := conf.strings[key]; present {
		conf.rwlock.RUnlock()
		return val
	}

	conf.rwlock.RUnlock()
	conf.SetString(key, def)
	return conf.GetString(key, def)
}

func (conf *GoNetConfig) Keys() []string {
	conf.rwlock.RLock()
	defer conf.rwlock.RUnlock()

	totalLen := len(conf.ints) + len(conf.floats) + len(conf.strings) + len(conf.bools)
	ret := make([]string, totalLen)
	idx := 0

	for key, _ := range(conf.ints) {
		ret[idx] = key
		idx += 1
	}

	for key, _ := range(conf.floats) {
		ret[idx] = key
		idx += 1
	}

	for key, _ := range(conf.strings) {
		ret[idx] = key
		idx += 1
	}

	for key, _ := range(conf.bools) {
		ret[idx] = key
		idx += 1
	}

	sort.Strings(ret)
	return ret
}

func NewGoNetConfig() *GoNetConfig {
	return &GoNetConfig{
		rwlock:  new(sync.RWMutex),
		ints:    make(map[string]int),
		floats:  make(map[string]float32),
		strings: make(map[string]string),
		bools:   make(map[string]bool),
	}
}
