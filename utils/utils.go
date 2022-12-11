package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrNon200   = errors.New("received non 200 response code")
	ErrNotFile  = errors.New("content-type isn't right")
	ErrTooSmall = errors.New("image too small")
)

func DownloadFile(URL string) ([]byte, error) {
	// Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		fmt.Println(URL)
		fmt.Printf("%T", err)
		fmt.Println("hi1", err)
		return nil, err

	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, ErrNon200
	}

	// TODO: Check if we need this part
	contentType := response.Header.Get("content-type")
	if contentType == "" {
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
