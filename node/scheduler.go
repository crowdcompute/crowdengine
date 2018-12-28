package node

import (
	"fmt"
	"strings"
	"time"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/manager"
	"github.com/docker/docker/api/types"
)

// Running for ever, or until job's done
func WaitForJobToFinish(containerID string) bool {
	fmt.Println("start task status tracking")
	// TODO: Time has to be a const somewhere
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking if job's done...")
			if !containerRunning(containerID) {
				return true
			}
		}
	}
}

func containerRunning(containerID string) bool {
	cjson, err := manager.GetInstance().InspectContainer(containerID)
	if err != nil {
		fmt.Println("Error inspecting container. ID : \n", containerID)
		return false
	}
	// If at least one is running then state that I am busy
	if cjson.State.Running {
		return true
	}
	return false
}

// Running for ever, or until node dies
func PruneImages(quit <-chan struct{}) {
	fmt.Println("start prunning images")
	// TODO: Time has to be a const somewhere
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking if there are images to be removed...")
			RemoveImages()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func RemoveImages() {
	summaries, err := manager.GetInstance().ListImages(types.ImageListOptions{All: true})
	common.CheckErr(err, "[ListImages] Failed to List images")
	now := time.Now().Unix()
	for _, img := range summaries {
		image := database.ImageLvlDB{}
		imgID := strings.Replace(img.ID, "sha256:", "", -1)

		i, err := database.GetDB().Model(image).Get([]byte(imgID))
		image, ok := i.(database.ImageLvlDB)
		if !ok {
			continue
		}
		// If the image was found
		if err == nil {
			if time.Unix(image.CreatedTime, 0).Add(common.TenDays).Unix() <= now {
				fmt.Println("Removing image: ", img.ID)
				manager.GetInstance().RemoveImage(img.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
			}
		}
	}
}
