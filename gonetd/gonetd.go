package main

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"log"
)

func main() {
	ctx := &daemon.Context{
		PidFileName: "gonetd-pid",
		PidFilePerm: 0644,
		LogFileName: "gonetd-log",
		LogFilePerm: 0640,
		WorkDir: "./",
		Umask: 027,
		Args: []string { "[gonetd]" },
	}

	fmt.Println(ctx)

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
	
	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error: ", err)
	}
}
