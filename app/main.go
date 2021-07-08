package main

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/shovanmaity/grpc-task/app/interceptor"
	"github.com/shovanmaity/grpc-task/app/server"
	"github.com/shovanmaity/grpc-task/app/store"
	"github.com/shovanmaity/grpc-task/app/types"
	generated "github.com/shovanmaity/grpc-task/gen/go"
)

func main() {
	serverPort := types.ServerPort()
	gatewayPort := types.GatewayPort()

	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		logrus.Fatalf("Failed to listen on port %s, error : %s", serverPort, err)
	}

	db := store.New()
	inc := interceptor.New(db)

	pingSRV := server.NewPingServer()
	loginSRV := server.NewLoginServer(db)
	profileSRV := server.NewProfileServer(db)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(inc.UnaryServerInterceptor()),
	)

	generated.RegisterPingServiceServer(s, pingSRV)
	generated.RegisterLoginServiceServer(s, loginSRV)
	generated.RegisterProfileServiceServer(s, profileSRV)

	logrus.Infof("Starting gRPC server on %s port", serverPort)
	go func() {
		logrus.Fatal(s.Serve(listener))
	}()

	conn, err := grpc.DialContext(context.Background(),
		":"+serverPort, grpc.WithBlock(), grpc.WithInsecure(),
	)
	if err != nil {
		logrus.Fatalf("Failed to dial , error : %s", err)
	}

	smux := runtime.NewServeMux()
	if err := generated.RegisterPingServiceHandler(context.Background(), smux, conn); err != nil {
		logrus.Fatalf("Failed to register ping gateway, error : %s", err)
	}
	if err := generated.RegisterLoginServiceHandler(context.Background(), smux, conn); err != nil {
		logrus.Fatalf("Failed to register login gateway, error : %s", err)
	}
	if err := generated.RegisterProfileServiceHandler(context.Background(), smux, conn); err != nil {
		logrus.Fatalf("Failed to register profile gateway, error : %s", err)
	}

	gwServer := &http.Server{
		Addr:    ":" + gatewayPort,
		Handler: smux,
	}

	logrus.Infof("Serving gRPC gateway on %s port", gatewayPort)
	logrus.Fatal(gwServer.ListenAndServe())
}
