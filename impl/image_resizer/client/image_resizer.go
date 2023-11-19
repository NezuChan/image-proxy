package client

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	"github.com/davidbyttow/govips/v2/vips"
	resty "github.com/go-resty/resty/v2"
)

type ImageResizerClient interface {
	DecryptImage(imageHex string) (string, error)
	GetOriginImage(url string) ([]byte, error)
	ResizeImage(image []byte, width, height int) ([]byte, error)
}

type imageResizerClient struct {
	httpClient resty.Client
	crypter    cipher.BlockMode
	iv         []byte
}

func (i imageResizerClient) ResizeImage(image []byte, width, height int) (result []byte, err error) {
img, err := vips.NewImageFromBuffer(image)
    defer img.Close()
    if err != nil {
        err = errors.New(fmt.Sprintf("unable to read image: %v", err.Error()))
        return
    }

    // Get original size of image
    originalWidth := float64(img.Width())
    originalHeight := float64(img.Height())
	scaleWidth := float64(width) / originalWidth
    scaleHeight := float64(height) / originalHeight

	scaleFactor := math.Min(scaleWidth, scaleHeight)

    if originalWidth == originalHeight {
		if err = img.Resize(scaleFactor, vips.KernelNearest); err != nil {
			return
		}
	} else {
		if height < width {
			width = height
		}

		if height > int(originalHeight) {
			height = int(originalHeight)
		}
		
    	// Crop the image to the specified dimensions while maintaining the aspect ratio
    	if err = img.SmartCrop(width, height, vips.InterestingCentre); err != nil {
        	return
    	}
	}


    options := vips.NewJpegExportParams()
    options.Quality = 100
    result, _, err = img.ExportJpeg(options)
    if err != nil {
        err = errors.New(fmt.Sprintf("unable to export image: %v", err))
        return
    }

    return
}

func (i imageResizerClient) GetOriginImage(url string) (result []byte, err error) {
	response, err := i.httpClient.R().Get(url)
	if err != nil {
		return
	}

	result = response.Body()
	return
}

func (i imageResizerClient) DecryptImage(imageHex string) (result string, err error) {
	image, err := hex.DecodeString(imageHex)
	if err != nil {
		err = errors.New("couldn't decode image hex")
		return
	}

	decryptedImage := make([]byte, len(image))
	i.crypter.CryptBlocks(decryptedImage, image)
	i.crypter.(interface{ SetIV([]byte) }).SetIV(i.iv) // Reset the IV

	result = string(decryptedImage[:len(decryptedImage)-int(decryptedImage[len(decryptedImage)-1])])
	return
}

func NewImageResizerClient(rawIV string, rawKey string) ImageResizerClient {
	key := []byte(rawKey)
	iv := []byte(rawIV)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("couldn't initialize new cipher: %v", err))
	}
	return &imageResizerClient{
		httpClient: *resty.New(),
		crypter:    cipher.NewCBCDecrypter(block, iv),
		iv:         iv,
	}
}
