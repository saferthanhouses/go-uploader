package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Endpoint string
var SpaceName string
var AccessKey string
var SecKey string
var AppSecret []byte
var Ssl bool

func init(){
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	Endpoint = os.Getenv("ENDPOINT")
	SpaceName = os.Getenv("SPACE_NAME")
	AccessKey = os.Getenv("SPACES_KEY")
	SecKey = os.Getenv("SPACES_SECRET")
	AppSecret = []byte(os.Getenv("APP_SECRET"))
	Ssl = true
}