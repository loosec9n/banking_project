package main

import (
	"fmt"
	"log"
	"net/http"
	"simplebank/controller"
)

func main() {
	fmt.Println("SimpleBank")

	//handler fucntion
	http.HandleFunc("/", controller.HelloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
