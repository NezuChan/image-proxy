package domain

type ImageResizerUsecase interface {
	GetMaxResolution() (int, int)
	ResolveImage(imageHex string, width, height int) ([]byte, error)
}

type ImageResizerHTTPResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
