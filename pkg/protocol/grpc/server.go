package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"github.com/neerajrush/go-grpc-http-rest/pkg/logger"
	"github.com/neerajrush/go-grpc-http-rest/pkg/api/v1"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Info("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Info("starting gRPC server...")
	return server.Serve(listen)
}
