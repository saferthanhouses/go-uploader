package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"learnabout-filemanager/upload-client"
	"learnabout-filemanager/utils"
	"log"
	"net/http"
	"strings"
)


type Upload struct {
	Url string
}

// http.HandlerFunc wraps the handler function in a struct, supplies this to the mux and calls it within the ServeHTTP method
var UploadHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	// Not a POST request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method Not Allowed")
		return
	}

	file, header, err := r.FormFile("data")
	if err != nil || file == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error accessing the file")
		return
	}

	defer file.Close()

	contentType, err := utils.GetFileContentType(file)
	content := strings.Split(contentType, "/")
	contentBase, contentSuffix := content[0], content[1]

	if err != nil || contentBase != "image" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No Image Found")
		log.Printf("%v", err)
		return
	}

	// Generate ID for the file
	id := xid.New()
	_, err = upload_client.Upload(id.String(), contentType, contentSuffix, file, header.Size)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error uploading the file")
	}

	url := utils.GetUploadUrl(id.String(), contentSuffix)

	upload := Upload {
		Url: url,
	}

	body, err := json.Marshal(&upload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
})
