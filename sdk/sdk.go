package gorigelsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("failed when caching image with statuscode = %v", resp.StatusCode)
		return "", err
	}
	var cacheImageResponse CacheImageResponse
	err = json.NewDecoder(resp.Body).Decode(&cacheImageResponse)
	if err != nil {
		return "", err
	}

	signedQueryString = SignQueryString(s.key, s.salt, fmt.Sprintf("img/%s", cacheImageResponse.Signature), "", expiry)
	pathURL = fmt.Sprintf("%s/img/%s?%s", s.baseURL, cacheImageResponse.Signature, signedQueryString)
	return pathURL, nil
}

func (s *SDK) BatchedCacheImage(proxyParamsSlice []ProxyParams, expiry int64) ([]CacheImageResponse, error) {
	signedQueryString := SignQueryString(s.key, s.salt, "batched-headsup", "", expiry)
	pathURL := fmt.Sprintf("%s/batched-headsup?%s", s.baseURL, signedQueryString)

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(proxyParamsSlice)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(pathURL, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("failed when caching image with statuscode = %v", resp.StatusCode)
		return nil, err
	}
	var cacheImageResponse []CacheImageResponse
	err = json.NewDecoder(resp.Body).Decode(&cacheImageResponse)

	// Create ShortURL
	for i := 0; i < len(cacheImageResponse); i++ {
		signedQueryString = SignQueryString(s.key, s.salt, fmt.Sprintf("img/%s", cacheImageResponse[i].Signature), "", expiry)
		cacheImageResponse[i].ShortURL = fmt.Sprintf("%s/img/%s?%s", s.baseURL, cacheImageResponse[i].Signature, signedQueryString)

	}
	return cacheImageResponse, err
}

func (sdk *SDK) TryShortURL(imageURL string, options *Options, expiry int64) string {
	pathURL, err := sdk.CacheImage(imageURL, options, expiry)
	if err != nil {
		return sdk.ProxyImage(imageURL, options, expiry)
	}

	return pathURL
}
