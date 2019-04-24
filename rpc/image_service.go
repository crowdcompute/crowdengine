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

	"github.com/crowdcompute/crowdengine/manager"
	"github.com/docker/docker/api/types"
)

// ImageService used to register Docker image functionality
// over jsonrpc
type ImageService struct{}

// NewImageService returns a new ImageService
func NewImageService() *ImageService {
	return &ImageService{}
}

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
	uploadPath, ok := r.Context().Value(common.ContextKeyUploadPath).(string)

	result, err := manager.GetInstance().LoadImage(uploadPath +"/uploads/" + filename)
	if err != nil {
		return result, err
	}
	return result, nil
}
