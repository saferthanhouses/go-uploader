package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go"
	"github.com/rs/xid"
	"log"
	"net/http"
	"os"
)


var endpoint string = "ams3.digitaloceanspaces.com"
var spaceName string = "learnabout-dev" // Space names must be globally unique

func main(){
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	/*
	The Spaces API aims to be interoperable with Amazon's AWS S3 API. In most cases,
	when using a client library, setting the "endpoint" or "base" URL to ${REGION}.digitaloceanspaces.com
	and generating a Spaces key to replace your AWS IAM key will allow you to use Spaces in place of S3.
 	*/

	accessKey := os.Getenv("SPACES_KEY")
	secKey := os.Getenv("SPACES_SECRET")
	ssl := true

	//Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)

	if err != nil {
		log.Fatal(err)
	}

	h := http.NewServeMux()

	h.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// this is implicitly setting the status header to 200
		//fmt.Fprintf(w, "Hello %s", r.Header.Get("Content-Type"))

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Not a post request")
			return
		} else {
			file, header, err := r.FormFile("data")
			if file == nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Error accessing the file")
				return
			}

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Error accessing the file")
				log.Printf("%v", err)
				return
			}

			// Cors stuff
			r.Header.Set("Access-Control-Allow-Origin", "http://localhost:8080")

			// Generate ID for the file
			id := xid.New()

			n, err := client.PutObject(spaceName, id.String(), file, header.Size, minio.PutObjectOptions{})

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Error uploading the file")
			}

			fmt.Fprintf(w, "File %s uploaded successfully. %d bytes", id, n)
		}
	})

	serverErr := http.ListenAndServe(":8081", h)
	log.Fatal(serverErr)
}