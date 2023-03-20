package usecase

import (
	"github.com/nezuchan/image-proxy/domain"
	"github.com/nezuchan/image-proxy/impl/image_resizer/client"
)

type imageResizerUsecase struct {
	client    client.ImageResizerClient
	maxWidth  int
	maxHeight int
}

func (i imageResizerUsecase) ResolveImage(imageHex string, width, height int) (result []byte, err error) {
	url, err := i.client.DecryptImage(imageHex)
	if err != nil {
		return
	}

	resp, err := i.client.GetOriginImage(url)
	if err != nil {
		return
	}

	result, err = i.client.ResizeImage(resp, width, height)
	if err != nil {
		return
	}

	return
}

func (i imageResizerUsecase) GetMaxResolution() (int, int) {
	return i.maxWidth, i.maxHeight
}

func NewImageResizerUsecase(client client.ImageResizerClient, maxWidth, maxHeight int) domain.ImageResizerUsecase {
	return &imageResizerUsecase{
		client:    client,
		maxWidth:  maxWidth,
		maxHeight: maxHeight,
	}
}
