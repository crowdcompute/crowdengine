package rpc

import (
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
	filename, fileHandler := getFileFromRequest(w, r)
	defer fileHandler.Close()

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
	fmt.Fprint(w, hexHash)
}

func getFileFromRequest(w http.ResponseWriter, r *http.Request) (string, multipart.File) {
	// Get the file from the http request
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*1024) // 500 Mb
	r.ParseMultipartForm(32 << 20)                          // 33 Mb memory
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, "Unable to upload file. Error: ", err, file)
		return "", nil
	}
	return handler.Filename, file
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
	log.Println("The hash is: ", hexHash)
	image := &database.ImageAccount{Signature: hexSignature, Path: path, CreatedTime: time.Now().Unix()}
	database.GetDB().Model(image).Put([]byte(hash))
	return hexHash
}
