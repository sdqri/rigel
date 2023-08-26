package gorigelsdk

import (
	"fmt"
	"net/url"
)

type ProxyParams struct {
	Img     string  `json:"img"`
	Options Options `json:"options"`
}

type Options struct {
	Height         *int
	Width          *int
	AreaHeight     *int
	AreaWidth      *int
	Top            *int
	Left           *int
	Quality        *int
	Compression    *int
	Zoom           *int
	Crop           *bool
	Enlarge        *bool
	Embed          *bool
	Flip           *bool
	Flop           *bool
	Force          *bool
	NoAutoRotate   *bool
	NoProfile      *bool
	Interlace      *bool
	StripMetadata  *bool
	Trim           *bool
	Lossless       *bool
	Extend         *string
	Rotate         *string
	Background     *string
	Gravity        IGravity
	Watermark      *string
	WatermarkImage *string
	Type           IImageType
	Interpolator   *string
	Interpretation *string
	GaussianBlur   *string
	Sharpen        *string
	Threshold      *float64
	Gamma          *float64
	Brightness     *float64
	Contrast       *float64
	OutputICC      *string
	InputICC       *string
	Palette        *bool
}

func NewOptions(source *Options) *Options {
	return &Options{
		Height:         source.Height,
		Width:          source.Width,
		AreaHeight:     source.AreaHeight,
		AreaWidth:      source.AreaWidth,
		Top:            source.Top,
		Left:           source.Left,
		Quality:        source.Quality,
		Compression:    source.Compression,
		Zoom:           source.Zoom,
		Crop:           source.Crop,
		Enlarge:        source.Enlarge,
		Embed:          source.Embed,
		Flip:           source.Flip,
		Flop:           source.Flop,
		Force:          source.Force,
		NoAutoRotate:   source.NoAutoRotate,
		NoProfile:      source.NoProfile,
		Interlace:      source.Interlace,
		StripMetadata:  source.StripMetadata,
		Trim:           source.Trim,
		Lossless:       source.Lossless,
		Extend:         source.Extend,
		Rotate:         source.Rotate,
		Background:     source.Background,
		Gravity:        source.Gravity,
		Watermark:      source.Watermark,
		WatermarkImage: source.WatermarkImage,
		Type:           source.Type,
		Interpolator:   source.Interpolator,
		Interpretation: source.Interpretation,
		GaussianBlur:   source.GaussianBlur,
		Sharpen:        source.Sharpen,
		Threshold:      source.Threshold,
		Gamma:          source.Gamma,
		Brightness:     source.Brightness,
		Contrast:       source.Contrast,
		OutputICC:      source.OutputICC,
		InputICC:       source.InputICC,
		Palette:        source.Palette,
	}
}

func (o *Options) QueryString() string {
	v := url.Values{}
	if o.Height != nil {
		v.Set("height", fmt.Sprintf("%v", *o.Height))
	}
	if o.Width != nil {
		v.Set("width", fmt.Sprintf("%v", *o.Width))
	}
	if o.AreaHeight != nil {
		v.Set("areaheight", fmt.Sprintf("%v", *o.AreaHeight))
	}
	if o.AreaWidth != nil {
		v.Set("areawidth", fmt.Sprintf("%v", *o.AreaWidth))
	}
	if o.Top != nil {
		v.Set("top", fmt.Sprintf("%v", *o.Top))
	}
	if o.Left != nil {
		v.Set("left", fmt.Sprintf("%v", *o.Left))
	}
	if o.Quality != nil {
		v.Set("quality", fmt.Sprintf("%v", *o.Quality))
	}
	if o.Compression != nil {
		v.Set("compression", fmt.Sprintf("%v", *o.Compression))
	}
	if o.Zoom != nil {
		v.Set("zoom", fmt.Sprintf("%v", *o.Zoom))
	}
	if o.Crop != nil {
		v.Set("crop", fmt.Sprintf("%v", *o.Crop))
	}
	if o.Enlarge != nil {
		v.Set("enlarge", fmt.Sprintf("%v", *o.Enlarge))
	}
	if o.Embed != nil {
		v.Set("embed", fmt.Sprintf("%v", *o.Embed))
	}
	if o.Flip != nil {
		v.Set("flip", fmt.Sprintf("%v", *o.Flip))
	}
	if o.Flop != nil {
		v.Set("flop", fmt.Sprintf("%v", *o.Flop))
	}
	if o.Force != nil {
		v.Set("force", fmt.Sprintf("%v", *o.Force))
	}
	if o.NoAutoRotate != nil {
		v.Set("noautorotate", fmt.Sprintf("%v", *o.NoAutoRotate))
	}
	if o.NoProfile != nil {
		v.Set("noprofile", fmt.Sprintf("%v", *o.NoProfile))
	}
	if o.Interlace != nil {
		v.Set("interlace", fmt.Sprintf("%v", *o.Interlace))
	}
	if o.StripMetadata != nil {
		v.Set("stripmetadata", fmt.Sprintf("%v", *o.StripMetadata))
	}
	if o.Trim != nil {
		v.Set("trim", fmt.Sprintf("%v", *o.Trim))
	}
	if o.Lossless != nil {
		v.Set("lossless", fmt.Sprintf("%v", *o.Lossless))
	}
	if o.Extend != nil {
		v.Set("extend", *o.Extend)
	}
	if o.Rotate != nil {
		v.Set("rotate", *o.Rotate)
	}
	if o.Background != nil {
		v.Set("background", *o.Background)
	}
	if o.Gravity != nil {
		v.Set("gravity", o.Gravity.String())
	}
	if o.Watermark != nil {
		v.Set("watermark", *o.Watermark)
	}
	if o.WatermarkImage != nil {
		v.Set("watermarkimage", *o.WatermarkImage)
	}
	if o.Type != nil {
		v.Set("type", o.Type.String())
	}
	if o.Interpolator != nil {
		v.Set("interpolator", *o.Interpolator)
	}
	if o.Interpretation != nil {
		v.Set("interpretation", *o.Interpretation)
	}
	if o.GaussianBlur != nil {
		v.Set("gaussianblur", *o.GaussianBlur)
	}
	if o.Sharpen != nil {
		v.Set("sharpen", *o.Sharpen)
	}
	if o.Threshold != nil {
		v.Set("threshold", fmt.Sprintf("%v", *o.Threshold))
	}
	if o.Gamma != nil {
		v.Set("gamma", fmt.Sprintf("%v", *o.Gamma))
	}
	if o.Brightness != nil {
		v.Set("brightness", fmt.Sprintf("%v", *o.Brightness))
	}
	if o.Contrast != nil {
		v.Set("contrast", fmt.Sprintf("%v", *o.Contrast))
	}
	if o.OutputICC != nil {
		v.Set("outputicc", fmt.Sprintf("%v", *o.OutputICC))
	}
	if o.InputICC != nil {
		v.Set("inputicc", fmt.Sprintf("%v", *o.InputICC))
	}
	if o.Palette != nil {
		v.Set("palette", fmt.Sprintf("%v", *o.Palette))
	}

	escapedQuery := v.Encode()

	unescapedQuery, _ := url.QueryUnescape(escapedQuery)
	return unescapedQuery
}
