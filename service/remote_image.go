package service

import (
	"encoding/json"
	"fmt"

	"github.com/h2non/bimg"
)

type RemoteImage struct {
	ImageRequest `json:"image_request"`
	Data         *[]byte `json:"data"`
}

type RemoteImageOption func(*RemoteImage)

func WithURL(url string) RemoteImageOption {
	return func(ri *RemoteImage) {
		ri.URL = url
	}
}

func WithbimgOptions(options bimg.Options) RemoteImageOption {
	return func(ri *RemoteImage) {
		ri.Options = options
	}
}

func WithImageRequest(imageRequest *ImageRequest) RemoteImageOption {
	return func(ri *RemoteImage) {
		ri.ImageRequest = *imageRequest
	}
}

func NewRemoteImage(opts ...RemoteImageOption) *RemoteImage {
	remoteImage := RemoteImage{}
	for _, opt := range opts {
		opt(&remoteImage)
	}
	return &remoteImage
}

func (remoteImage *RemoteImage) String() string {
	return fmt.Sprintf("[%v]%s", len(*remoteImage.Data), remoteImage.URL)
	return remoteImage.String()
}

func (remoteImage *RemoteImage) Process() (err error) {
	bimgImage := bimg.NewImage(*remoteImage.Data)
	processedData, err := bimgImage.Process(remoteImage.Options)
	if err != nil {
		return
	}

	remoteImage.Data = &processedData
	return nil
}

func (remoteImage *RemoteImage) Type() string {
	bimgImage := bimg.NewImage(*remoteImage.Data)
	return bimgImage.Type()
}

func (remoteImage *RemoteImage) getValue() (value string, err error) {
	valueData, err := json.Marshal(remoteImage)
	if err != nil {
		return
	}
	value = string(valueData)
	return
}

func (remoteImage *RemoteImage) GetPair(prefix string) (key string, value string, err error) {
	key = remoteImage.GetKey(prefix)
	value, err = remoteImage.getValue()
	return
}

func (remoteImage *RemoteImage) ParseValue(value string) error {
	return json.Unmarshal([]byte([]byte(value)), remoteImage)
}

// func FromCacheable(cacheable adapters.Cacheable) (remoteImage RemoteImage, err error) {
// 	cacheableJson, err := json.Marshal(cacheable)
// 	if err != nil {
// 		return
// 	}
// 	err = json.Unmarshal(cacheableJson, &remoteImage)
// 	if err != nil {
// 		return
// 	}
// 	return
// }
