package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

var (
	serverPort = os.Getenv("PORT")
)

func main() {
	router := httprouter.New()
	router.POST("/register", RegisterController)
	router.POST("/login", LoginController)
	router.GET("/resource", ListResourceController)
	router.POST("/resource", CreateResourceController)
	router.POST("/machine", CreateMachineController)
	router.POST("/lease", CreateLeaseController)
	router.GET("/lease", ListLeasesController)
	router.DELETE("/lease", DeleteLeaseController)

	// add middleware
	serve := alice.New().Then(router)

	log.Fatal(http.ListenAndServe(":"+serverPort, serve))
}
