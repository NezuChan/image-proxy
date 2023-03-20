package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/discord/lilliput"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	godotenv.Load()
	app := fiber.New()

	key := []byte(os.Getenv("KEY"))
	iv := []byte(os.Getenv("IV"))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	app.Get("/:size/:image", func(c *fiber.Ctx) error {
		size := strings.Split(c.Params("size", "512x512"), "x")
		x, err := strconv.Atoi(size[0])
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message":    "Invalid size",
				"statusCode": 400,
			})
		}

		y, err := strconv.Atoi(size[1])
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message":    "Invalid size",
				"statusCode": 400,
			})
		}

		width, err := strconv.Atoi(os.Getenv("MAX_WIDTH"))
		height, err := strconv.Atoi(os.Getenv("MAX_HEIGHT"))

		if x > width || y > height {
			return c.Status(400).JSON(fiber.Map{
				"message":    "Image too large",
				"statusCode": 400,
			})
		}

		image, err := hex.DecodeString(fmt.Sprintf(c.Params("image")))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message":    "image decode fail",
				"statusCode": 500,
			})
		}

		decrypted := make([]byte, len(image))
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(decrypted, image)

		resp, err := http.Get(string(decrypted[:len(decrypted)-int(decrypted[len(decrypted)-1])]))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		imageData, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		decoder, _ := lilliput.NewDecoder(imageData)
		header, _ := decoder.Header()

		ops := lilliput.NewImageOps(8192)
		defer ops.Close()

		outputImg := make([]byte, 8*x*y)

		if x == 0 {
			x = header.Width()
		}

		if y == 0 {
			y = header.Height()
		}

		resizeMethod := lilliput.ImageOpsResize
		if x == header.Width() && y == header.Height() {
			resizeMethod = lilliput.ImageOpsNoResize
		}

		opts := &lilliput.ImageOptions{
			FileType:             ".jpeg",
			Width:                x,
			Height:               y,
			ResizeMethod:         resizeMethod,
			NormalizeOrientation: true,
			EncodeOptions:        map[int]int{lilliput.JpegQuality: 100},
		}

		outputImg, _ = ops.Transform(decoder, opts, outputImg)
		c.Set("Content-Type", "image/jpeg")
		return c.Send(outputImg)
	})

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	app.Listen(fmt.Sprintf("%s:%s", host, port))
}
