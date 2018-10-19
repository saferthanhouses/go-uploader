package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go"
	"github.com/rs/xid"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type Upload struct {
	Url string
}

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
	appSecret := []byte(os.Getenv("APP_SECRET"))

	//Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)

	if err != nil {
		log.Fatal(err)
	}

	// Create a client for our authentication service
	//db := prisma.New(&prisma.PrismaOptions{
	//	Endpoint: "http://localhost:4466/",
	//})
	//ctx := context.Background()

	h := http.NewServeMux()

	h.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		authenticated, err := authenticate(r, appSecret)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server Error")
			return
		}

		if authenticated == false {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Not Authorized")
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Not a post request")
			return
		}

		file, header, err := r.FormFile("data")

		if err != nil || file == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error accessing the file")
			return
		}

		defer file.Close()

		contentType, err := getFileContentType(file)
		content := strings.Split(contentType, "/")
		contentBase, contentSuffix := content[0], content[1]

		if err != nil || contentBase != "image" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "No Image Found")
			log.Printf("%v", err)
			return
		}

		// Cors stuff
		r.Header.Set("Access-Control-Allow-Origin", "http://localhost:8080")

		// Generate ID for the file
		id := xid.New()

		// Set extra content-policy keys
		userMetaData := map[string]string{"x-amz-acl": "public-read"}
		_, err = client.PutObject(spaceName, id.String() + "." + contentSuffix, file, header.Size, minio.PutObjectOptions{
			ContentType: contentType,
			UserMetadata: userMetaData,
		})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error uploading the file")
		}

		w.WriteHeader(http.StatusOK)

		urlComponents := []string {
			endpoint,
			spaceName,
			id.String() + "." + contentSuffix,
		}

		upload := Upload {
			Url: "http://" + strings.Join(urlComponents, "/"),
		}

		body, err := json.Marshal(&upload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server Error")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	serverErr := http.ListenAndServe(":8081", h)
	log.Fatal(serverErr)
}

func authenticate(r *http.Request, appSecret []byte) (bool, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return false, nil
	}

	authorization = strings.Replace(authorization, "Bearer ", "", 1)

	parser := jwt.Parser{
		SkipClaimsValidation: true,
	}

	token, err := parser.Parse(authorization, func(t *jwt.Token) (interface{}, error){

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return appSecret, nil
	})

	if err != nil {
		fmt.Printf("error in parsing token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["userId"])
		return true, nil
	} else {
		fmt.Println(err)
		return false, nil
	}
}

// Brazenly Stolen from here: https://golangcode.com/get-the-content-type-of-file/
func getFileContentType(file multipart.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)

	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}