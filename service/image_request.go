package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/h2non/bimg"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sdqri/rigel/utils"
)

var (
	ErrInvalidSignature       error = errors.New("invalid signature")
	ErrResourceClaimNotExists error = errors.New("resource claim does not exists")
	ErrResourceClaimType            = errors.New("resource claim cannot cast into string type")
)

type AlgKey struct {
	Alg    jwa.KeyAlgorithm
	PubKey jwk.Key
}

type ImageRequest struct {
	URL       string       `json:"url"`
	Options   bimg.Options `json:"options"`
	Signature string       `json:"-"`
}

func (ir *ImageRequest) String() string {
	return ir.URL
}

func (ir *ImageRequest) GetKey(prefix string) (key string) {
	key = fmt.Sprintf("%s:%s", strings.TrimPrefix(prefix, "/"), ir.Signature)
	return
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

func ParseToken(algKey AlgKey, queryToken string, debug bool) (*ImageRequest, error) {
	token, err := jwt.ParseString(queryToken, jwt.WithKey(algKey.Alg, algKey.PubKey), jwt.WithVerify(!debug))
	if err != nil {
		return nil, err
	}
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
	sha1 := sha1.New()
	sha1.Write([]byte(queryToken))
	shasum := sha1.Sum(nil)
	ImageRequest := ImageRequest{
		Signature: hex.EncodeToString(shasum),
	}
	fmt.Println(ImageRequest)
	err = json.Unmarshal([]byte(res), &ImageRequest)
	if err != nil {
		return nil, err
	}

	return &ImageRequest, nil
}

func NewImageRequest(signature string) *ImageRequest {
	ImageRequest := ImageRequest{
		Signature: signature,
	}
	return &ImageRequest
}
