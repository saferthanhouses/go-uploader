package main

import (
	"learnabout-filemanager/middleware"
	"learnabout-filemanager/routes"
	"log"
	"net/http"
)


func main(){
	h := http.NewServeMux()
	h.Handle("/upload", middleware.Cors(middleware.Authentication(routes.UploadHandler)))

	serverErr := http.ListenAndServe(":8081", h)
	log.Fatal(serverErr)
}