package main

import (
	"github.com/joho/godotenv"
	"github.com/minio/minio-go"
	"github.com/valyala/fasthttp"
	"log"
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

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)

	if err != nil {
		log.Fatal(err)
	}

	// pass plain function to fasthttp
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/upload":
			uploadHandler(ctx, client)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}

	fasthttp.ListenAndServe(":8081", requestHandler)
}

func uploadHandler(ctx *fasthttp.RequestCtx, client *minio.Client) {


	n, err := client.FPutObject(spaceName, objectName, filePath, minio.PutObjectOptions{ContentType:contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}