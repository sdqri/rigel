package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/h2non/bimg"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sdqri/rigel/utils"
)

var (
	ErrInvalidSignature       error = errors.New("invalid signature")
	ErrResourceClaimNotExists error = errors.New("resource claim does not exists")
	ErrResourceClaimType            = errors.New("resource claim cannot cast into string type")
)

type ImageRequest struct {
	URL     string       `json:"url"`
	Options bimg.Options `json:"options"`
}

func (ir *ImageRequest) String() string {
	return ir.URL
}

func (ir *ImageRequest) GetKey(prefix string) (key string) {
	return fmt.Sprintf("%s:%s", strings.TrimPrefix(prefix, "/"), ir.URL)
}

func (ir *ImageRequest) Download() (remoteImage *RemoteImage, err error) {
	data, err := utils.DownloadFile(ir.URL)
	if err != nil {
		return
	}
	remoteImage = NewRemoteImage(WithImageRequest(ir))
	remoteImage.Data = &data
	return
}

func ParseToken(queryToken string) (*ImageRequest, error) {
	token, err := jwt.ParseString(queryToken, jwt.WithVerify(false)) //TODO: Add verification

	// Getting res key
	resJSON, ok := token.Get("res")
	if !ok {
		return nil, ErrResourceClaimNotExists
	}

	res, ok := resJSON.(string)
	if !ok {
		return nil, ErrResourceClaimType
	}

	// Unmarshaling res into ImageRequest
	ImageRequest := ImageRequest{}
	err = json.Unmarshal([]byte(res), &ImageRequest)
	if err != nil {
		return nil, err
	}

	return &ImageRequest, nil
}
