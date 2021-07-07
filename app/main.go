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
	generated "github.com/shovanmaity/grpc-task/gen/go"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatalf("Failed to listen on port 8080, error : %s", err)
	}
	db := store.NewStore()

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

	logrus.Println("Starting gRPC server on 8080 port")
	go func() {
		logrus.Fatal(s.Serve(listener))
	}()

	conn, err := grpc.DialContext(context.Background(),
		":8080", grpc.WithBlock(), grpc.WithInsecure(),
	)
	if err != nil {
		logrus.Fatalf("Failed to dial server, error : %s", err)
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
		Addr:    ":8090",
		Handler: smux,
	}

	logrus.Info("Serving gRPC-Gateway on 8090 port")
	logrus.Fatal(gwServer.ListenAndServe())
}
