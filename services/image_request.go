package services

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/h2non/bimg"
	"github.com/sdqri/rigel/adapters"
	"github.com/sdqri/rigel/utils"
)

// TODO:
// - GaussianBlur
// - Sharpen
// - Background
// - Watermark
// - WatermarkImage
type IROptions struct {
	Height         int     `json:"Height,omitempty" form:"height"`
	Width          int     `json:"Width,omitempty" form:"width"`
	AreaHeight     int     `json:"AreaHeight,omitempty" form:"area_height"`
	AreaWidth      int     `json:"AreaWidth,omitempty" form:"area_width"`
	Top            int     `json:"Top,omitempty" form:"top"`
	Left           int     `json:"Left,omitempty" form:"left"`
	Quality        int     `json:"Quality,omitempty" form:"quality"`
	Compression    int     `json:"Compression,omitempty" form:"compression"`
	Zoom           int     `json:"Zoom,omitempty" form:"zoom"`
	Crop           bool    `json:"Crop,omitempty" form:"crop"`
	Enlarge        bool    `json:"Enlarge,omitempty" form:"enlarge"`
	Embed          bool    `json:"Embed,omitempty" form:"embed"`
	Flip           bool    `json:"Flip,omitempty" form:"flip"`
	Flop           bool    `json:"Flop,omitempty" form:"flop"`
	Force          bool    `json:"Force,omitempty" form:"force"`
	NoAutoRotate   bool    `json:"NoAutoRotate,omitempty" form:"no_auto_rotate"`
	NoProfile      bool    `json:"NoProfile,omitempty" form:"no_profile"`
	Interlace      bool    `json:"Interlace,omitempty" form:"interlace"`
	StripMetadata  bool    `json:"StripMetadata,omitempty" form:"strip_metadata"`
	Trim           bool    `json:"Trim,omitempty" form:"trim"`
	Lossless       bool    `json:"Lossless,omitempty" form:"lossless"`
	Extend         int     `json:"Extend,omitempty" form:"extend"`                 //Extend
	Rotate         int     `json:"Rotate,omitempty" form:"rotate"`                 //Angle
	Gravity        int     `json:"Gravity,omitempty" form:"gravity"`               //Gravity
	Type           int     `json:"Type,omitempty" form:"type"`                     //ImageType
	Interpolator   int     `json:"Interpolator,omitempty" form:"interpolator"`     //Interpolator
	Interpretation int     `json:"Interpretation,omitempty" form:"interpretation"` //Interpretation
	Threshold      float64 `json:"Threshold,omitempty" form:"threshold"`
	Gamma          float64 `json:"Gamma,omitempty" form:"gamma"`
	Brightness     float64 `json:"Brightness,omitempty" form:"brightness"`
	Contrast       float64 `json:"Contrast,omitempty" form:"contrast"`
	OutputICC      string  `json:"OutputICC,omitempty" form:"output_icc"`
	InputICC       string  `json:"InputICC,omitempty" form:"input_icc"`
	Palette        bool    `json:"Palette,omitempty" form:"palette"`
	Speed          int     `json:"Speed,omitempty" form:"speed"`
}

func (iro *IROptions) ToBimgOptions() (options bimg.Options, err error) {
	iroJson, err := json.Marshal(iro)
	if err != nil {
		return
	}
	err = json.Unmarshal(iroJson, &options)
	if err != nil {
		return
	}
	return
}

type ImageRequest struct {
	URL       string    `json:"url"`     // image url
	Options   IROptions `json:"options"` // transformation options
	signature string    //sha1sum
}

func NewImageRequest(rawURL string, options IROptions) (imageRequest *ImageRequest, err error) {
	_, err = url.Parse(rawURL)
	if err != nil {
		return
	}
	imageRequest = &ImageRequest{
		URL:     rawURL,
		Options: options,
	}
	return
}

func NewQueryableImageRequest(signature string) *ImageRequest {
	return &ImageRequest{
		signature: signature,
	}
}

func (ir *ImageRequest) String() string {
	key, _ := ir.GetKey()
	return fmt.Sprintf("ImageRequest(rawURL=%s, key=%s)", ir.URL, key)
}

func (ir *ImageRequest) SHA1Sum() (sha1sum string, err error) {
	irJson, err := json.Marshal(ir)
	if err != nil {
		return
	}
	sha1sum = utils.SHA1Sum(string(irJson))
	return
}

func (ir *ImageRequest) GetKey() (key string, err error) {
	var signature string
	if ir.signature != "" {
		signature = ir.signature
	} else {
		signature, err = ir.SHA1Sum()
		if err != nil {
			return
		}
	}
	key = fmt.Sprintf("image_request:%s", signature)
	return
}

func (ir *ImageRequest) GetValue() (value string, err error) {
	valueData, err := json.Marshal(ir)
	if err != nil {
		return
	}
	value = string(valueData)
	return
}

func (ir *ImageRequest) GetPair() (kv adapters.Pair[string], err error) {
	key, err := ir.GetKey()
	if err != nil {
		return
	}
	value, err := ir.GetValue()
	if err != nil {
		return
	}
	kv = adapters.Pair[string]{
		Key:   key,
		Value: value,
	}
	return
}

func (ir *ImageRequest) ParseValue(value string) error {
	err := json.Unmarshal([]byte(value), ir)
	if err != nil {
		return err
	}
	return nil
}

func (ir *ImageRequest) GetQueryableRemoteImage() (remoteImage *RemoteImage, err error) {
	key, err := ir.GetKey()
	if err != nil {
		return
	}
	remoteImage = &RemoteImage{
		SHA1Sum: key,
	}
	return
}

func (ir *ImageRequest) DownloadImage() (remoteImage *RemoteImage, err error) {
	// Getting key for
	key, err := ir.GetKey()
	if err != nil {
		return
	}
	data, err := utils.DownloadFile(ir.URL)
	if err != nil {
		return
	}
	remoteImage, err = NewRemoteImage(key, &data)
	if err != nil {
		return
	}
	return
}

func (ir *ImageRequest) GetRemoteImage() (remoteImage *RemoteImage, err error) {
	// download & process
	remoteImage, err = ir.DownloadImage()
	if err != nil {
		return
	}
	err = remoteImage.Process(ir.Options)
	return
}
