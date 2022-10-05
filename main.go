package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	db "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/pb"
	"simplebank/utils"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGRPCServer(config, store)

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
