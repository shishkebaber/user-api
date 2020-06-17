package main

import (
	"github.com/nicholasjackson/env"
	"github.com/shishkebaber/user-api/server"
	"os"
	"os/signal"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Server bind address")
var postgresBindAddress = env.String("PG_BIND_ADDRESS", false, "", "Server bind address")

func main() {
	env.Parse()

	s := server.NewServer(bindAddress, postgresBindAddress)

	exitChan := make(chan struct{}, 1)
	mainExit := make(chan os.Signal, 1)

	go s.Start(exitChan)

	signal.Notify(mainExit, os.Interrupt)
	signal.Notify(mainExit, os.Kill)
	<-mainExit
	exitChan <- struct{}{}
}
