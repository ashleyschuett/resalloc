package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	role       = os.Getenv("ROLE") // master / slave
	serverPort = os.Getenv("PORT")
)

func main() {
	router := httprouter.New()

	router.POST("/register", RegisterController)
	router.POST("/login", LoginController)

	// add middleware
	serve := alice.New().Then(router)

	log.Fatal(http.ListenAndServe(":"+serverPort, serve))
}
