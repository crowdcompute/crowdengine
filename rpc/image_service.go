package rpc

import (
	"context"

	"github.com/crowdcompute/crowdengine/manager"
	"github.com/docker/docker/api/types"
)

// ImageService used to register Docker image functionality
// over jsonrpc
type ImageService struct{}

// List list all the available images
func (s *ImageService) List(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	images, err := manager.GetInstance().ListImages(options)
	if err != nil {
		return images, err
	}
	return images, nil
}

// BuildFromDockerfile builds image from docker file
func (s *ImageService) BuildFromDockerfile(ctx context.Context) Result {
	return Result{"not implemented"}
}

// Load loads a tar
func (s *ImageService) Load(ctx context.Context, filename string) (string, error) {
	result, err := manager.GetInstance().LoadImage(filename)
	if err != nil {
		return result, err
	}
	return result, nil
}
