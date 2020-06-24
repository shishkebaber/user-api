package server

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/shishkebaber/user-api/data"
	"github.com/shishkebaber/user-api/handlers"
	protos "github.com/shishkebaber/user-api/protos/user"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type Server struct {
	http.Server
	Logger       *log.Logger
	UserHandlers *handlers.Users
	grpcS        *protos.UserGrpcServer
}

func NewServer() *Server {
	l := log.New()

	// Getting ENVs
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
	var grpcPort = os.Getenv("GRPC_PORT")

	// Initialising postgresDB, and making migrations
	postgresDB := data.NewPgDb(l, data.GenerateURL(pgUser, pgPassword, pgHost, pgPort, pgName, l))
	data.CreateTable(postgresDB, l) //TODO: use third party package for migrations

	// Initialising validation
	v := data.NewValidation()

	// Channel to receive updates on users
	userUpdateChan := make(chan *protos.UserData, 1)

	// Init gRPC server, for notifying subscribed clients
	grpcS := protos.NewUserGrpcServer(grpcPort, l, userUpdateChan)

	// Init handlers
	usersHandlers := handlers.NewUsersHandler(l, postgresDB, v, userUpdateChan)
	router := handlers.InitHandlers(usersHandlers)

	httpS := http.Server{
		Addr:         bindAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{httpS, l, usersHandlers, grpcS}
}

// Runs the http and grpc servers. Inputs: exitChan - channel to stop all servers
func (s *Server) Start(exitChan chan struct{}) {
	s.Logger.Info("Server starting")
	// Starting grpc server
	go s.grpcS.StartGrpc(exitChan)
	err := s.ListenAndServe()
	if err != nil {
		s.Logger.Error("Error during starting the server", err)
		os.Exit(1)
	}

	sig := <-exitChan
	s.Logger.Println("Got signal: ", sig)

	s.UserHandlers.Db.(*data.UserPostgresDb).Pool.Close()
	s.Logger.Info("Server stopping")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(ctx)
}
