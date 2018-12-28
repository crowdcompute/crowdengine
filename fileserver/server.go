package fileserver

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// FileServer allows files upload over HTTP
type FileServer struct {
}

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString generates a random number of string
func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	for i := 0; i < strlen; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}

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
	filename := RandomString(30) + filepath.Ext(handler.Filename)
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
