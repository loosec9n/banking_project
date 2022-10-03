package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"simplebank/api"
	"simplebank/controller"
	db "simplebank/db/sqlc"
	"simplebank/utils"

	_ "github.com/lib/pq"
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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create an server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start serer:", err)
	}

	//handler fucntion
	http.HandleFunc("/", controller.HelloHandler)
	log.Println(http.ListenAndServe(":8080", nil))
}
