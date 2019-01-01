// Copyright 2018 The crowdcompute:crowdengine Authors
// This file is part of the crowdcompute:crowdengine library.
//
// The crowdcompute:crowdengine library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The crowdcompute:crowdengine library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the crowdcompute:crowdengine library. If not, see <http://www.gnu.org/licenses/>.

package fileserver

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/crowdcompute/crowdengine/common"
)

// FileServer allows files upload over HTTP
type FileServer struct {
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
