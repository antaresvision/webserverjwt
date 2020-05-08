package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SayHello(writer http.ResponseWriter, request *http.Request) {
	//extract variables from request (both path and querystring variables -- NO POST variables here)
	vars := mux.Vars(request)
	name := vars["name"]
	if name == "" {
		name = "World"
	}
	response := fmt.Sprintf("Hello, %s!\n", name)
	writer.Write([]byte(response))
}
