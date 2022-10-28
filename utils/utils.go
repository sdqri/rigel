package utils

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

var (
	ErrNon200 = errors.New("received non 200 response code")
	ErrNotFile = errors.New("content-type isn't right")
	ErrTooSmall = errors.New("image too small")
)

func DownloadFile(URL string) ([]byte, error) {
	// Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, ErrNon200
	}

	contentType := response.Header.Get("content-type")
	if contentType == "" || !strings.HasPrefix(contentType, "image") {
		return nil, ErrNotFile
	}

	file, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if len(file) < 10 {
		return nil, ErrTooSmall
	}

	return file, nil
}
