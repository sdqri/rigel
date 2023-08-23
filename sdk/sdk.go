package gorigelsdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sdqri/rigel/controllers"
)

type SDK struct {
	baseURL string
	key     string
	salt    string
}

func NewSDK(baseURL string, key string, salt string) *SDK {
	return &SDK{
		baseURL: baseURL,
		key:     key,
		salt:    salt,
	}
}

func (s *SDK) ProxyImage(imageURL string, options *Options, expiry int64) string {
	queryString := ""
	if options != nil && options.QueryString() != "" {
		queryString = SerializeMapToQueryString(map[string]string{"img": imageURL}) + "&" + options.QueryString()
	} else {
		queryString = SerializeMapToQueryString(map[string]string{"img": imageURL})
	}
	signedQueryString := SignQueryString(s.key, s.salt, "proxy", queryString, expiry)
	pathURL := fmt.Sprintf("%s/proxy?%s", s.baseURL, signedQueryString)
	return pathURL
}

func (s *SDK) CacheImage(imageURL string, options *Options, expiry int64) (string, error) {
	queryString := ""
	if options != nil && options.QueryString() != "" {
		queryString = SerializeMapToQueryString(map[string]string{"img": imageURL}) + "&" + options.QueryString()
	} else {
		queryString = SerializeMapToQueryString(map[string]string{"img": imageURL})
	}
	signedQueryString := SignQueryString(s.key, s.salt, "headsup", queryString, expiry)
	pathURL := fmt.Sprintf("%s/headsup?%s", s.baseURL, signedQueryString)
	resp, err := http.Post(pathURL, "", nil)
	if err != nil {
		err := fmt.Errorf("error while trying to cacheImage with imageURL=%s: %v\n", imageURL, err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("failed when caching image with statuscode = %v", resp.StatusCode)
		return "", err
	}
	var cacheImageResponse controllers.CacheImageResponse
	err = json.NewDecoder(resp.Body).Decode(&cacheImageResponse)
	if err != nil {
		return "", err
	}
	return cacheImageResponse.Signature, nil
}

func (sdk *SDK) TryShortURL(imageURL string, options *Options, expiry int64) string {
	signature, err := sdk.CacheImage(imageURL, options, expiry)
	if err != nil {
		return sdk.ProxyImage(imageURL, options, expiry)
	}

	signedQueryString := SignQueryString(sdk.key, sdk.salt, fmt.Sprintf("img/%s", signature), "", expiry)
	pathURL := fmt.Sprintf("%s/img/%s?%s", sdk.baseURL, signature, signedQueryString)
	return pathURL
}
