package main

import (
	"log"
	"net/http"
	"os"
	"simplebank/controller"
)

func main() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.Println("SimpleBank")

	//handler fucntion
	http.HandleFunc("/", controller.HelloHandler)
	log.Println(http.ListenAndServe(":8080", nil))
}
