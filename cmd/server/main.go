package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/erpaher/grpc-crud/pkg/app"
	"github.com/erpaher/grpc-crud/pkg/store"
	api "github.com/erpaher/grpc-crud/proto/serviceexampleapi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	config := app.LoadConfig()
	store := store.Database(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.StoreUser, config.StorePassword, config.StoreHost, config.StorePort, config.StoreDatabase))
	proxy := app.New(store)

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	api.RegisterServiceExampleServiceServer(s, proxy)
	// Serve gRPC Server
	grpcAddr := fmt.Sprintf("%s:%d", config.GRPCHost, config.GRPCPort)
	log.Printf("Serving gRPC on %s\n", grpcAddr)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		grpcAddr,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = api.RegisterServiceExampleServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwAddr := fmt.Sprintf("%s:%d", config.GatewayHost, config.GatewayPort)
	gwServer := &http.Server{
		Addr:    gwAddr,
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http://%s\n", gwAddr)
	log.Fatalln(gwServer.ListenAndServe())
}
