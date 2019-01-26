package rpc

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/crowdcompute/crowdengine/common"
)

// ServeHTTP accepts file uploads multipart/form-data
// should return an id which represents the uploaded file
// should be able to register a clientID with a list of uploaded files
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the file from the http request
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*1024) // 500 Mb
	r.ParseMultipartForm(32 << 20)                          // 33 Mb memory
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, "Unable to upload file. Error: ", err, file)
		return
	}
	defer file.Close()
	// Save the file
	filename := common.RandomString(30) + filepath.Ext(handler.Filename)
	f, err := os.OpenFile("./uploads/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	// Return a response to the requester
	fmt.Fprint(w, filename)
}
