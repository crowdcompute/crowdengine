package node

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/crowdcompute/crowdengine/p2p"
	peer "github.com/libp2p/go-libp2p-peer"
)

type ImageManagerAPI struct {
	host   *p2p.Host
	images map[string][]byte // image hash -> signature
}

func NewImageManagerAPI(h *p2p.Host) *ImageManagerAPI {
	return &ImageManagerAPI{
		host: h,
		// TODO: NOT SURE IF THIS IS A GOOD IDEA
		images: make(map[string][]byte),
	}
}

// API call to push an image to the remote peer
func (api *ImageManagerAPI) PushImage(ctx context.Context, nodePID string, imageFilePath string) (string, error) {

	file, err := os.Open(imageFilePath)
	common.CheckErr(err, "[PushImage] Couldn't open file.")
	defer file.Close()

	peerid, err := peer.IDB58Decode(nodePID)
	common.CheckErr(err, "[PushImage] Couldn't IDB58Decode nodePID.")

	fileInfo, err := file.Stat()
	common.CheckErr(err, "[PushImage] Couldn't find stats.")

	// Hash image
	hash := crypto.HashImagePath(imageFilePath)
	hashString := hex.EncodeToString(hash)
	signature := api.images[hashString]
	signatureString := hex.EncodeToString(signature)

	// Starting a new stream to send a file
	api.host.SetConsistentStream(peerid)

	// TODO: all those numbers should go as constants
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	// TODO: Not sure what number to give here. Need to see the range
	filledSignature := fillString(signatureString, 150)
	// TODO: Not sure what number to give here. Need to see the range
	filledHash := fillString(hashString, 100)
	fmt.Println("Sending filename and filesize!")
	fmt.Println("fileSize: ", fileSize)
	fmt.Println("fileName: ", fileName)
	fmt.Println("filledSignature: ", filledSignature)
	fmt.Println("filledHash: ", filledHash)

	api.host.UploadChunk([]byte(fileSize))
	api.host.UploadChunk([]byte(fileName))
	api.host.UploadChunk([]byte(filledSignature))
	api.host.UploadChunk([]byte(filledHash))
	sendBuffer := make([]byte, common.FileChunk)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		api.host.UploadChunk(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return <-api.host.ImageIDchan, nil
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func (api *ImageManagerAPI) RunImage(ctx context.Context, nodePID string, imageID string) string {
	toNodeID, _ := peer.IDB58Decode(nodePID)
	api.host.RunImage(toNodeID, imageID)

	// Check if there are any pending requests to run
	containerID := <-api.host.ContainerID
	fmt.Println("Result running the job: ", containerID)
	return containerID
}

func (api *ImageManagerAPI) InspectContainer(ctx context.Context, nodePID string, containerID string) (string, error) {
	toNodeID, _ := peer.IDB58Decode(nodePID)
	api.host.CreateSendInspectRequest(toNodeID, containerID)
	fmt.Println("Result running the job: ")
	return <-api.host.InspectChan, nil
}

// Uploading an image to the current node
// TODO: Upload an image to the supernode instead of passing the file path
func (api *ImageManagerAPI) UploadImage(ctx context.Context, imageFilePath string, privateKey string) (string, error) {
	privByte, _ := hex.DecodeString(privateKey)
	priv, err := crypto.RestorePrivateKey(privByte)

	// Hash image
	// TODO: bytes will be received straight away, not from a path
	hash := crypto.HashImagePath(imageFilePath)
	// content, _ := ioutil.ReadFile(imageFilePath)
	sign, err := priv.Sign(hash)
	common.CheckErr(err, "[UploadImage] Failed to sign image.")
	api.images[hex.EncodeToString(hash)] = sign
	return "", nil
}

// Getting the list of images specific to the publicKey
func (api *ImageManagerAPI) ListImages(ctx context.Context, nodePID string, publicKey string) (string, error) {
	toNodeID, _ := peer.IDB58Decode(nodePID)
	api.host.CreateAndSendListRequest(toNodeID, publicKey)
	return <-api.host.ListChan, nil
}

// ************************************************************************************** //
// 								The old way											      //
/*
// API call to push and run an image to the remote peer
func (api *ImageManagerAPI) RunImageOnNode(ctx context.Context, nodePID string, imageFilePath string) (string, error) {
	fmt.Printf("<%s> Running imageID %s to node %s\n", api.protocolTask.P2PHost.ID().String(), imageFilePath, nodePID)

	// In the case we send the image via a POST request to the node
	result, err := api.UploadAndRunImage(ctx, api.getNodeAddr(nodePID), imageFilePath)

	// This is in the case the image is on a Registry,
	// the imageID or key needs to be passed to the Node for it to load it.

	// nodeID, _ := peer.IDB58Decode(nodePID)
	// api.node.nodeExecImage(nodeID, imageFilePath)

	return result, err
}

func (api *ImageManagerAPI) getNodeAddr(nodePID string) []ma.Multiaddr {
	peerid, err := peer.IDB58Decode(nodePID)
	common.CheckErr(err, "[getNodeAddr] Couldn't IDB58Decode nodePID.")
	pi := api.protocolTask.P2PHost.Peerstore().PeerInfo(peerid)
	return pi.Addrs
}

func (api *ImageManagerAPI) UploadAndRunImage(ctx context.Context, nodeAddr []ma.Multiaddr, imageFilePath string) (string, error) {
	fmt.Println("I got this image:", imageFilePath)

	// TODO: Need to change this logic and make it more generic
	// Why always take the first nodeAddr[0]?
	// Is IP always on the 3rd position?
	uri := nodeAddr[0].String()
	components := strings.Split(uri, "/")
	url := "http://" + components[2] + ":" + strconv.Itoa(api.node.uploadAddrPort) + "/upload"
	fmt.Println(url)

	request, err := newfileUploadRequest(url, "file", imageFilePath)
	common.CheckErr(err, "[UploadAndRunImage] Couldn't upload file.")

	client := &http.Client{}
	resp, err := client.Do(request)
	common.CheckErr(err, "[UploadAndRunImage] Couldn't do http request.")

	var bodyContent []byte
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	resp.Body.Read(bodyContent)
	fmt.Println(bodyContent)
	resp.Body.Close()
	return string(bodyContent), nil
}
*/
//***************************************************************************************//
//*******************************// Helper functions //**********************************//
//***************************************************************************************//
/*
// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}
*/
