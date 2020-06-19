package server

import (
	"context"
	"github.com/joho/godotenv"
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

func NewServer() *Server {
	l := log.New()

	err := godotenv.Load()
	if err != nil {
		l.Fatalf("Error getting env: %v", err)
	} else {
		l.Println("Getting the env values")
	}

	var bindAddress = os.Getenv("BIND_ADDRESS")
	var pgUser = os.Getenv("POSTGRES_USER")
	var pgPassword = os.Getenv("POSTGRES_PASSWORD")
	var pgName = os.Getenv("POSTGRES_DB")
	var pgHost = os.Getenv("POSTGRES_HOST")
	var pgPort = os.Getenv("POSTGRES_PORT")

	postgresDB := data.NewPgDb(l, data.GenerateURL(pgUser, pgPassword, pgHost, pgPort, pgName, l))
	v := data.NewValidation()

	data.CreateTable(postgresDB, l) //TODO: use third party package for migrations

	usersHandlers := handlers.NewUsersHandler(l, postgresDB, v)

	router := handlers.InitHandlers(usersHandlers)

	httpS := http.Server{
		Addr:         bindAddress,
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
