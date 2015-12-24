//
//  Configuration Keys
//
package gonet

import (
	"errors"
	"fmt"
	"sync"
)

const (
	IPV4_PACKET_CHANNEL_SIZE = "ipv4_packet_channel_size"
)

type GoNetConfig struct {
	rwlock *sync.RWMutex
	ints map[string]int
}

func (conf *GoNetConfig) GetInt(key string) (int, error) {
	conf.rwlock.RLock()
	defer conf.rwlock.RUnlock()
	
	if val, present := conf.ints[key]; present {
		return val, nil
	}
	return -1, errors.New(fmt.Sprintf("Integer key '%s' not present in configuration.", key))
}

func NewGoNetConfig() (*GoNetConfig) {
	return &GoNetConfig{
		rwlock: new(sync.RWMutex),
		ints: make(map[string]int),
	}
}
