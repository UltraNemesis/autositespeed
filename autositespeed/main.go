// main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/UltraNemesis/autositespeed"
)

var conf autositespeed.Configuration

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	autositespeed.LoadConfig([]string{"./conf", "../conf"}, "config", &conf)

	server := autositespeed.NewServer(conf)

	fmt.Println("Starting AutoSiteSpeed Services...")

	go server.Start()

	<-sigs
}
