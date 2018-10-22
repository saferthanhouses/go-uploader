package upload_client

import (
	"github.com/minio/minio-go"
	"io"
	"learnabout-filemanager/config"
	"log"
)

var Client *minio.Client

func init(){
	var err error
	//Initiate a client using DigitalOcean Spaces.
	Client, err = minio.New(config.Endpoint, config.AccessKey, config.SecKey, config.Ssl)

	if err != nil {
		log.Fatal(err)
	}

}

func Upload(id string, contentType string, contentSuffix string, file io.Reader, size int64) (int64, error){
	// Set extra content-policy keys for the uploaded file
	userMetaData := map[string]string{"x-amz-acl": "public-read"}

	n, err := Client.PutObject(
		config.SpaceName, id + "." + contentSuffix, file, size, minio.PutObjectOptions{
			ContentType: contentType,
			UserMetadata: userMetaData,
		})

	if err != nil {
		log.Printf("attempting to upload: %v", err)
	}

	return n, err
}