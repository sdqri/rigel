package services

import (
	"encoding/json"
	"fmt"

	"github.com/h2non/bimg"
	"github.com/sdqri/rigel/adapters"
)

type RemoteImage struct {
	// ImageRequest `json:"image_request"`
	SHA1Sum string  `json:"-"`
	Data    *[]byte `json:"data"`
}

func NewRemoteImage(sha1sum string, data *[]byte) (remoteImage *RemoteImage, err error) {
	remoteImage = &RemoteImage{
		SHA1Sum: sha1sum,
		Data:    data,
	}
	return
}

func (ri *RemoteImage) String() string {
	key, _ := ri.GetKey()
	return fmt.Sprintf("RemoteImage(key=%s)", key)
}

func (ri *RemoteImage) GetKey() (key string, err error) {
	key = fmt.Sprintf("remote_image:%s", ri.SHA1Sum)
	return
}

func (ri *RemoteImage) GetValue() (value string, err error) {
	valueData, err := json.Marshal(ri)
	if err != nil {
		return
	}
	value = string(valueData)
	return
}

func (ri *RemoteImage) GetPair() (kv adapters.Pair[string], err error) {
	key, err := ri.GetKey()
	if err != nil {
		return
	}
	value, err := ri.GetValue()
	if err != nil {
		return
	}
	kv = adapters.Pair[string]{
		Key:   key,
		Value: value,
	}
	return
}

func (ri *RemoteImage) ParseValue(value string) error {
	err := json.Unmarshal([]byte(value), ri)
	if err != nil {
		return err
	}
	return nil
}

func (ri *RemoteImage) Type() string {
	bimgImage := bimg.NewImage(*ri.Data)
	return bimgImage.Type()
}

func (ri *RemoteImage) FileName() string {
	return fmt.Sprintf("%s.%s", ri.SHA1Sum, ri.Type())
}

func (ri *RemoteImage) ContentType() string {
	return fmt.Sprintf("image/%s", ri.Type())
}

func (ri *RemoteImage) Process(options IROptions) (err error) {
	bimgImage := bimg.NewImage(*ri.Data)
	bimgOptions, err := options.ToBimgOptions()
	if err != nil {
		return
	}
	processedData, err := bimgImage.Process(bimgOptions)
	if err != nil {
		return
	}
	ri.Data = &processedData
	return nil
}
