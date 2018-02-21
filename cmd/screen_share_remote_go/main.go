package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rootkiwi/screen_share_remote_go/cli"
	"github.com/rootkiwi/screen_share_remote_go/conf"
	"github.com/rootkiwi/screen_share_remote_go/remote"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		switch arg := args[0]; arg {
		case "-h", "--help", "help":
			cli.PrintUsage()
		case "genconf":
			cli.GenConf()
		case "noconf":
			conf := cli.NoConf()
			start(conf)
		default:
			conf, err := conf.ParseConfigFile(arg)
			if err != nil {
				log.Fatalln(err)
			}
			start(conf)
		}
	} else {
		cli.PrintUsage()
	}
}

func start(conf *conf.Config) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM)
	done := make(chan struct{}, 1)
	go remote.Listen(conf, done)
	log.Println("screen_share_remote started")
	<-signals
	done <- struct{}{}
	select {
	case <-done:
	case <-time.After(time.Second * 2): // waiting max 2 sec for cleanup
	}
}
