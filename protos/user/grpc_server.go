package protos

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
	"sync"
)

type UserGrpcServer struct {
	port       string
	logger     *logrus.Logger
	UserUpdate chan *UserData
	subs       map[UserService_SubscribeOnUsersUpdatesServer]struct{}
}

func (uGrpc *UserGrpcServer) removeSubscriber(stream UserService_SubscribeOnUsersUpdatesServer) {
	if _, ok := uGrpc.subs[stream]; ok {
		delete(uGrpc.subs, stream)
	}
}
func (uGrpc *UserGrpcServer) SubscribeOnUsersUpdates(stream UserService_SubscribeOnUsersUpdatesServer) error {
	for {
		_, err := stream.Recv()
		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			uGrpc.logger.Info("Client has closed connection.")
			uGrpc.removeSubscriber(stream)
			break
		}

		if err != nil {
			uGrpc.logger.Error("Unable to read from client.", "error", err)
			uGrpc.removeSubscriber(stream)
			break
		}

		uGrpc.logger.Info("Handle client subscribe request.")
		uGrpc.subs[stream] = struct{}{}
	}

	return nil
}

func (uGrpc *UserGrpcServer) mustEmbedUnimplementedUserServiceServer() {
	panic("implement me")
}

func NewUserGrpcServer(port string, l *logrus.Logger, userUpdtChan chan *UserData) *UserGrpcServer {
	return &UserGrpcServer{port, l, userUpdtChan, make(map[UserService_SubscribeOnUsersUpdatesServer]struct{})}
}

// Starting grpc server and handle function for users updates processing
func (uGrpc *UserGrpcServer) StartGrpc(stop chan struct{}) {
	uGrpc.logger.Info("Starting GRPC service at port: ", uGrpc.port)
	gs := grpc.NewServer()
	RegisterUserServiceServer(gs, uGrpc)
	reflection.Register(gs)
	go uGrpc.handleUsersUpdates()
	l, err := net.Listen("tcp", uGrpc.port)
	if err != nil {
		uGrpc.logger.Error("Unable to create listener ", "error: ", err)
		return
	}

	// listen for requests
	gs.Serve(l)
	<-stop
	gs.GracefulStop()
}

func (uGrpc *UserGrpcServer) handleUsersUpdates() {
	for user := range uGrpc.UserUpdate {
		uGrpc.logger.Info("UpdatedUserReceived: ", user.String())
		go uGrpc.broadcastUser(user)
	}
}

// Concurrent  iterating through subs
func (uGrpc *UserGrpcServer) broadcastUser(user *UserData) {
	done := make(chan struct{})
	wait := sync.WaitGroup{}
	uGrpc.logger.Info("Broadcasting user to ", len(uGrpc.subs), " subscribers. User Data:  ", user.String())
	for k, _ := range uGrpc.subs {
		wait.Add(1)

		go func(u *UserData, str UserService_SubscribeOnUsersUpdatesServer) {
			defer wait.Done()

			if _, ok := uGrpc.subs[str]; ok {
				err := str.Send(u)
				if err != nil {
					uGrpc.logger.Error("Unable to send updated user with ID: ", u.Id)
				}
			} else {
				uGrpc.logger.Error("Unable to send updated user, connection is closed")
			}
		}(user, k)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
