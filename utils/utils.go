package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

var (
	ErrNon200   = errors.New("received non 200 response code")
	ErrNotFile  = errors.New("content-type isn't right")
	ErrTooSmall = errors.New("image too small")
)

func DownloadFile(URL string) ([]byte, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Googlebot")

	client := http.Client{
		Jar: jar,
	}

	// Get the response bytes from the url
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
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
