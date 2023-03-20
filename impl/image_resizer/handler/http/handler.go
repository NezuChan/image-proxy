package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nezuchan/image-proxy/domain"
	"net/http"
	"strconv"
	"strings"
)

type httpMemberHandler struct {
	iUsecase domain.ImageResizerUsecase
}

func ConfigureImageResizerHandler(app *fiber.App, iUsecase domain.ImageResizerUsecase) {
	h := httpMemberHandler{
		iUsecase: iUsecase,
	}

	app.Add(http.MethodGet, "/:size/:image", h.GetImage)
}

func (h *httpMemberHandler) GetImage(c *fiber.Ctx) (err error) {
	size := strings.Split(c.Params("size", "512x512"), "x")

	width, err := strconv.Atoi(size[0])
	if err != nil {
		return c.Status(400).JSON(domain.ImageResizerHTTPResponse{
			StatusCode: 400,
			Message:    "invalid width size",
		})
	}

	height, err := strconv.Atoi(size[1])
	if err != nil {
		return c.Status(400).JSON(domain.ImageResizerHTTPResponse{
			StatusCode: 400,
			Message:    "invalid height size",
		})
	}

	maxWidth, maxHeight := h.iUsecase.GetMaxResolution()
	if width > maxWidth || height > maxHeight {
		return c.Status(400).JSON(domain.ImageResizerHTTPResponse{
			StatusCode: 400,
			Message:    "width or height size too large",
		})
	}

	image, err := h.iUsecase.ResolveImage(c.Params("image"), width, height)
	if err != nil {
		return c.Status(500).JSON(domain.ImageResizerHTTPResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	c.Set("Content-Type", "image/jpeg")
	return c.Send(image)
}
