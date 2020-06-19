package main

import (
	"github.com/shishkebaber/user-api/server"
	"os"
	"os/signal"
)

func main() {

	s := server.NewServer()

	exitChan := make(chan struct{}, 1)
	mainExit := make(chan os.Signal, 1)

	go s.Start(exitChan)

	signal.Notify(mainExit, os.Interrupt)
	signal.Notify(mainExit, os.Kill)
	<-mainExit
	exitChan <- struct{}{}
}
