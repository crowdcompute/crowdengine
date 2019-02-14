package rpc

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/crowdcompute/crowdengine/database"
)

// ServeHTTP accepts file uploads multipart/form-data
// should return an id which represents the uploaded file
// should be able to register a clientID with a list of uploaded files
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key, ok := r.Context().Value(common.ContextKeyPrivateKey).(*keystore.Key)
	if !ok {
		fmt.Fprintln(w, "There was an error getting the key from the context")
	}
	uploadPath, ok := r.Context().Value(common.ContextKeyUploadPath).(string)
	if !ok {
		fmt.Fprintln(w, "There was an error getting the upload path from the context")
	}

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

	filePath := uploadPath + "/uploads/" + filename

	// TODO: Why 0777 gets wrxr-xr-x
	const dirPerm = 0777
	if err := os.MkdirAll(filepath.Dir(filePath), dirPerm); err != nil {
		fmt.Fprint(w, err)
		return
	}
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	hash := crypto.HashFile(f)
	sign, err := key.KeyPair.Private.Sign(hash)
	hexHash := hex.EncodeToString(hash)
	storeImageToDB(hex.EncodeToString(hash), filePath, hex.EncodeToString(sign))

	// Return a response to the requester
	fmt.Fprint(w, hexHash)
}

// storeImageToDB stores the new image's data to our level DB
func storeImageToDB(hash, path, signature string) {
	image := &database.ImageAccount{Signature: signature, Path: path, CreatedTime: time.Now().Unix()}
	database.GetDB().Model(image).Put([]byte(hash))
}
