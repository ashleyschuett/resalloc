package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RegisterController handles account creation
func RegisterController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s RegisterModel
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&s)
	fmt.Println(s)
}

// LoginController handles authentication and token generation
func LoginController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
