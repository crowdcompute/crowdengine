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

package rpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/log"
	libcrypto "github.com/libp2p/go-libp2p-crypto"
)

func ServeFilesHTTP(ks *keystore.KeyStore, uploadPath string) http.HandlerFunc {
	return uploadAuthorization(ks, uploadPath, fileserve)
}

// UploadAuth authenticates a token and enriches the requests
// Authenticates a token and passes the request to the next handler
func uploadAuthorization(ks *keystore.KeyStore, uploadPath string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := getKeyForAccount(ks, r.Header)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), common.ContextKeyPrivateKey, key)
		ctx = context.WithValue(ctx, common.ContextKeyUploadPath, uploadPath)
		log.Printf("Token valid and account {%s} unlocked. ", key.Address)
		next(w, r.WithContext(ctx))
	}
}

// fileserve accepts file uploads multipart/form-data
// should return an id which represents the uploaded file
// should be able to register a clientID with a list of uploaded files
func fileserve(w http.ResponseWriter, r *http.Request) {
	key, ok := r.Context().Value(common.ContextKeyPrivateKey).(*keystore.Key)
	if !ok {
		fmt.Fprintln(w, "There was an error getting the key from the context")
	}
	uploadPath, ok := r.Context().Value(common.ContextKeyUploadPath).(string)
	if !ok {
		fmt.Fprintln(w, "There was an error getting the upload path from the context")
	}
	filename, fileHandler := getFileFromRequest(w, r)
	defer fileHandler.Close()
	if checkIfFileUploaded(fileHandler) {
		msg := fmt.Sprintf("File %s uploaded already", filename)
		log.Println(msg)
		fmt.Fprint(w, msg)
		return
	}
	fileHandler.Seek(0, 0)
	localFile, fullpath, err := createFile(filename, uploadPath)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	defer localFile.Close()
	_, err = io.Copy(localFile, fileHandler)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	// Rewind the file pointer to the beginning
	localFile.Seek(0, 0)
	log.Println("The file has been successgully uploaded, full path is: ", fullpath)
	hexHash := storeImageToDB(localFile, key.KeyPair.Private, fullpath)
	log.Println("The hash is: ", hexHash)
	fmt.Fprint(w, hexHash)
}

func getFileFromRequest(w http.ResponseWriter, r *http.Request) (string, multipart.File) {
	// Get the file from the http request
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*1024) // 500 Mb
	r.ParseMultipartForm(32 << 20)                          // 33 Mb memory
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error Here : ", err, file)
		fmt.Fprintln(w, "Unable to upload file. Error: ", err, file)
		return "", nil
	}
	return handler.Filename, file
}

func checkIfFileUploaded(f multipart.File) bool {
	hexHash := hex.EncodeToString(crypto.HashFile(f))
	_, err := database.GetImageAccountFromDB(hexHash)
	if err == database.ErrNotFound {
		return false
	} else if err != nil {
		log.Println("There was an error getting the image from DB.")
		return false
	}
	return true
}

func createFile(filename, path string) (*os.File, string, error) {
	randFilename := common.RandomString(30) + filepath.Ext(filename)
	fullpath := path + "/uploads/" + randFilename
	// TODO: Why 0777 gets wrxr-xr-x
	const dirPerm = 0777
	if err := os.MkdirAll(filepath.Dir(fullpath), dirPerm); err != nil {
		return nil, "", err
	}
	f, err := os.Create(fullpath)
	if err != nil {
		return nil, "", err
	}
	return f, fullpath, nil
}

// storeImageToDB stores the new image's data to our level DB
func storeImageToDB(f *os.File, priv libcrypto.PrivKey, path string) string {
	hash := crypto.HashFile(f)
	sign, err := priv.Sign(hash)
	common.FatalIfErr(err, "Couldn't sign with key")
	hexHash := hex.EncodeToString(hash)
	hexSignature := hex.EncodeToString(sign)
	image := &database.ImageAccount{Signature: hexSignature, Path: path, CreatedTime: time.Now().Unix()}
	database.GetDB().Model(image).Put([]byte(hexHash))
	return hexHash
}
