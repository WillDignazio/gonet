package main

import (
	"digitalbebop.net/gonet"
	"flag"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
)

func Init(config *gonet.GoNetConfig) {
	log.Println("Initializing...")
	_ = gonet.NewGoNetCtx(config)

	for idx, key := range config.Keys() {
		log.Println(idx,"\t", key)
	}
}

var configFile = flag.String("config", "/path/to/gonetd.conf", "Path to configuration file.")

func main() {
	flag.Parse()

	fmt.Println("ConfigFile: ", *configFile)
	config := ParseConfigFile(*configFile)

	logFile := config.GetString(LOG_PATH_KEY, "gonetd.log")
	logPerm := config.GetInt(LOG_PERM_KEY, 0600)
	pidFile := config.GetString(PID_PATH_KEY, "gonetd.pid")
	pidPerm := config.GetInt(PID_PERM_KEY, 0600)
	umask := config.GetInt(UMASK_KEY, 027)
	workDir := config.GetString(WORKDIR_PATH_KEY, "./")

	ctx := &daemon.Context{
		PidFileName: pidFile,
		PidFilePerm: os.FileMode(pidPerm),
		LogFileName: logFile,
		LogFilePerm: os.FileMode(logPerm),
		WorkDir:     workDir,
		Umask:       umask,
		Args:        flag.Args(),
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := ctx.Search()
		if err != nil {
			log.Fatalln("Failed to signal daemon: ", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := ctx.Reborn()
	if err != nil {
		log.Fatalln(err)
	}

	if d != nil {
		return
	}
	defer ctx.Release()

	log.Println("Starting gonetd.....")
	Init(config)

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error: ", err)
		return
	}
}
