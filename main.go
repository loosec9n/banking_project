package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	db "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/pb"
	"simplebank/utils"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println("SimpleBank")

	//configuration file
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("not able to load configurations ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("not able to connect to database :", err)
	}

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGRPCServer(config, store)

}

func runGatewayServer(config utils.Config, store *db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create an server: ", err)
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBankNowHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("not able to register the GRPC Gateway handler server - ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start the listener in GRPC server : ", err)
	}

	log.Printf("starting HTTP Gateway server in %s ", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("failed to start HTTP Gateway server : ", err)
	}
}
func runGRPCServer(config utils.Config, store *db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create an server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBankNowServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot start the listener in GRPC server : ", err)
	}

	log.Printf("starting GRPC server in %s ", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("failed to start GRPC server : ", err)
	}
}

// func runGinServer(config utils.Config, store *db.Store) {

// 	server, err := api.NewServer(config, store)
// 	if err != nil {
// 		log.Fatal("cannot create an server: ", err)
// 	}

// 	err = server.Start(config.HTTPServerAddress)
// 	if err != nil {
// 		log.Fatal("cannot start serer:", err)
// 	}

// 	//handler fucntion
// 	http.HandleFunc("/", controller.HelloHandler)
// 	log.Println(http.ListenAndServe(":8080", nil))
// }
