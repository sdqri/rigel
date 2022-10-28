package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"rigel/adapters"
	"rigel/utils"

	"github.com/h2non/bimg"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	ErrInvalidSignature error = errors.New("invalid signature")
	ErrResourceClaimNotExists error = errors.New("resource claim does not exists")
	ErrResourceClaimType = errors.New("resource claim cannot cast into string type")
)

type RemoteImage struct {
	URL string `json:"url"`
	Options bimg.Options `json:"options"`
	Data *[]byte `json:"data"`
	Prefix string `json:"prefix"`
	Downloaded bool `json:"downloaded"`
	Processed bool `json:"processed"`
	err error
}

type RemoteImageOption func(*RemoteImage)

func WithURL(url string) RemoteImageOption{
	return func(ri *RemoteImage) {
		ri.URL = url
	}
}

func WithbimgOptions(options bimg.Options) RemoteImageOption{
	return func(ri *RemoteImage) {
		ri.Options = options
	}
}

func WithPrefix(prefix string) RemoteImageOption{
	return func(ri *RemoteImage) {
		ri.Prefix = prefix
	}
}


func NewRemoteImage(opts ...RemoteImageOption) *RemoteImage{
	remoteImage := RemoteImage{}
	for _, opt := range opts {
		opt(&remoteImage)
	}
	return &remoteImage
}

func (remoteImage *RemoteImage) String() string {
	return remoteImage.GetKey()
}

func (remoteImage *RemoteImage) Assign(assignFrom adapters.Cacheable) error{
	afJson, err := json.Marshal(assignFrom)
	if err!=nil{
		return err
	}
	err = json.Unmarshal(afJson,&remoteImage)
	if err!=nil{
		return err
	}
	return nil
}

func (remoteImage *RemoteImage) Download() (err error) {
	if remoteImage.Downloaded == true{
		return remoteImage.err
	}

	if remoteImage.err != nil{
		return remoteImage.err
	}

	var data []byte
	data, err = utils.DownloadFile(remoteImage.URL)
	if err != nil {
		remoteImage.err = err
		return
	}
	remoteImage.Data = &data
	remoteImage.Downloaded = true
	return
}

func (remoteImage *RemoteImage) Redownload() (err error) {
	remoteImage.Downloaded = false
	remoteImage.err = nil
	return remoteImage.Download()
}


func (remoteImage *RemoteImage) Process() (err error) {
	if remoteImage.Processed == true {
		return nil
	}
	
	err = remoteImage.Download()
	if err!= nil{
		return
	}

	bimgImage := bimg.NewImage(*remoteImage.Data)
	processedData, err := bimgImage.Process(remoteImage.Options)
	if err != nil{
		return
	}

	remoteImage.Data = &processedData
	remoteImage.Processed = true
	return nil
}

func (remoteImage *RemoteImage) Reprocess() (err error) {
	remoteImage.Processed = false
	return remoteImage.Process()
}

func (remoteImage *RemoteImage) Type() (format string, err error) {
	err = remoteImage.Download()
	if err!= nil{
		return
	}
	bimgImage := bimg.NewImage(*remoteImage.Data)
	format = bimgImage.Type()
	return
}

func (remoteImage *RemoteImage) GetKey() (key string) {
	return fmt.Sprintf("%s:%s", strings.TrimPrefix(remoteImage.Prefix, "/"), remoteImage.URL)
}

func (remoteImage *RemoteImage) GetRedisPair() (key string, value string, err error) {
	key = remoteImage.GetKey()
	valueData, err := json.Marshal(remoteImage)
	if err != nil{
		return
	}
	value = string(valueData)
	return
}

func (remoteImage *RemoteImage) ParseValue(value string) error {
	return json.Unmarshal([]byte([]byte(value)), remoteImage)
}

func (remoteImage *RemoteImage) ParseQueryParams(queryParams string) error {
	token, err := jwt.ParseString(queryParams, jwt.WithVerify(false))

	resJSON, ok := token.Get("res")
	if !ok {
		return ErrResourceClaimNotExists
	}

	res, ok := resJSON.(string)
	if !ok {
		return ErrResourceClaimType
	}

	err = json.Unmarshal([]byte(res), &remoteImage)
	if err!= nil {
		return err
	}

	return nil
}
