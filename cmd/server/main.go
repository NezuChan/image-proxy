package main

import (
	"errors"
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/nezuchan/image-proxy/domain"
	"github.com/nezuchan/image-proxy/impl/image_resizer/client"
	i_http "github.com/nezuchan/image-proxy/impl/image_resizer/handler/http"
	i_uc "github.com/nezuchan/image-proxy/impl/image_resizer/usecase"
	"log"
	"os"
	"strconv"
)

func main() {
	// Initialize vips
	vips.LoggingSettings(func(domain string, level vips.LogLevel, msg string) {
		fmt.Println(domain, level, msg)
	}, vips.LogLevelInfo)

	vips.Startup(nil)
	defer vips.Shutdown()

	_ = godotenv.Load()
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "Fiber",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return ctx.Status(code).JSON(domain.ImageResizerHTTPResponse{
				StatusCode: code,
				Message:    err.Error(),
			})
		},
	})

	app.Use(recover.New())

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} [${latency}] ${status} - ${method} ${path}\n",
	}))

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	imageResizerClient := client.NewImageResizerClient(os.Getenv("IV"), os.Getenv("KEY"))
	maxWidth, err := strconv.Atoi(os.Getenv("MAX_WIDTH"))
	maxHeight, err := strconv.Atoi(os.Getenv("MAX_HEIGHT"))

	if err != nil {
		panic(err)
	}

	iUsecase := i_uc.NewImageResizerUsecase(imageResizerClient, maxWidth, maxHeight)
	i_http.ConfigureImageResizerHandler(app, iUsecase)
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", host, port)))
}
