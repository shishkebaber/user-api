package server

import (
	"context"
	"github.com/shishkebaber/user-api/data"
	"github.com/shishkebaber/user-api/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type Server struct {
	http.Server
	Logger       *log.Logger
	UserHandlers *handlers.Users
}

func NewServer(bindAddress *string, dbAddress *string) *Server {
	l := log.New()
	postgresDB := data.NewPgDb(l, *dbAddress)
	v := data.NewValidation()
	usersHandlers := handlers.NewUsersHandler(l, postgresDB, v)

	router := handlers.InitHandlers(usersHandlers)

	httpS := http.Server{
		Addr:         *bindAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
		IdleTimeout:  120 * time.Second,
	}

	return &Server{httpS, l, usersHandlers}
}

func (s *Server) Start(exitChan chan struct{}) {
	s.Logger.Info("Server starting")
	err := s.ListenAndServe()
	if err != nil {
		s.Logger.Error("Error during starting the server", err)
		os.Exit(1)
	}

	sig := <-exitChan
	log.Println("Got signal: ", sig)

	s.UserHandlers.Db.(*data.UserPostgresDb).Pool.Close()
	s.Logger.Info("Server stopping")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(ctx)
}
