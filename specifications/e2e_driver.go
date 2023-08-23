package specifications

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type E2EDriver struct {
	BaseURL string
	Client  *http.Client
}

func (driver *E2EDriver) Version() Result[string] {
	path, err := url.JoinPath(driver.BaseURL, "version")
	if err != nil {
		return Result[string]{Ok: "", Err: err}
	}
	resp, err := driver.Client.Get(path)
	if err != nil {
		return Result[string]{Ok: "", Err: err}
	}
	defer resp.Body.Close()

	var versionMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&versionMap)
	if err != nil {
		return Result[string]{Ok: "", Err: err}
	}
	version, ok := versionMap["version"]
	if !ok {
		return Result[string]{Ok: "", Err: err}
	}
	return Result[string]{Ok: version, Err: err}
}

func (driver *E2EDriver) Proxy(pathParams map[string]string) Result[*[]byte] {
	path, err := url.JoinPath(driver.BaseURL, "proxy")
	if err != nil {
		return Result[*[]byte]{Ok: nil, Err: err}
	}

	// Add path parameters
	path = buildURLWithParams(path, pathParams)

	resp, err := driver.Client.Get(path)
	if err != nil {
		return Result[*[]byte]{Ok: nil, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result[*[]byte]{Ok: nil, Err: err}
	}

	return Result[*[]byte]{Ok: &body, Err: err}
}

// func (driver *E2EDriver) HeadsUp() Result {
//
// }
//
// func (driver *E2EDriver) ValidateSignature() Result {
//
// }

func buildURLWithParams(baseURL string, params map[string]string) string {
	u, _ := url.Parse(baseURL)

	// Convert the map of path parameters to a slice of key-value pairs
	var pathSegments []string
	for key, value := range params {
		pathSegments = append(pathSegments, key, value)
	}

	// Join the path segments and set it as the URL's RawPath
	u.RawPath = strings.Join(pathSegments, "/")

	return u.String()
}
