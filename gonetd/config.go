package main

import (
	"digitalbebop.net/gonet"
	"fmt"
)

const (
	LOG_PATH_KEY     = "gonetd.log.path"
	LOG_PERM_KEY     = "gonetd.log.perm"
	PID_PATH_KEY     = "gonetd.pid.path"
	PID_PERM_KEY     = "gonetd.pid.perm"
	UMASK_KEY	 = "gonetd.umask"
	WORKDIR_PATH_KEY = "gonetd.work.path"
)

func ParseConfigFile(path string) *gonet.GoNetConfig {
	config := gonet.NewGoNetConfig()
	fmt.Println("XXX TODO: Parse Config File")
	// XXX TODO
	return config
}
