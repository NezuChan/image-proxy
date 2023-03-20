package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
	"image/jpeg"
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

	app.Use(logger.New())

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
			return c.Status(500).JSON(fiber.Map{
				"message":    "Unable to fetch origin image",
				"statusCode": 500,
			})
		}

		defer resp.Body.Close()

		imageData, _ := jpeg.Decode(resp.Body)

		outputImg := resize.Resize(uint(x), uint(y), imageData, resize.NearestNeighbor)
		buf := new(bytes.Buffer)
		options := &jpeg.Options{Quality: 100}
		err = jpeg.Encode(buf, outputImg, options)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message":    "Unable to encode image.",
				"statusCode": 500,
			})
		}
		output := buf.Bytes()

		c.Set("Content-Type", "image/jpeg")
		return c.Send(output)
	})

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	app.Listen(fmt.Sprintf("%s:%s", host, port))
}
