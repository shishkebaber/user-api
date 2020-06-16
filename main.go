package main

import (
	"context"
	"github.com/shishkebaber/user-api/data"
	"github.com/shishkebaber/user-api/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	log "github.com/sirupsen/logrus"
)

var bindAddress = env.String("BIND_ADRESS", false, ":9090", "Server bind address")

func main() {
	env.Parse()
	l := log.New()

	sMux := mux.NewRouter()

	postgresDB := data.NewPgDb(l)

	v := data.NewValidation()

	usersHandler := handlers.NewUsers(l, postgresDB, v)

	//Handlers
	getR := sMux.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/users", usersHandler.ListAll)

	s := http.Server{
		Addr:         *bindAddress,
		Handler:      sMux,
		ReadTimeout:  5 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Info("Server starting")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error during starting the server", err)
			os.Exit(1)
		}
	}()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, os.Kill)

	sig := <-exitChan
	log.Println("Got signal: ", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
