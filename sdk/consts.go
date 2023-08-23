package gorigelsdk

import (
	"fmt"
	"strconv"
)

type IImageType interface {
	GetImageType() string
	fmt.Stringer
}

type ImageType int

func (t ImageType) GetImageType() string {
	switch t {
	case JPEG:
		return "JPEG"
	case WEBP:
		return "WEBP"
	case PNG:
		return "PNG"
	case TIFF:
		return "TIFF"
	case GIF:
		return "GIF"
	case PDF:
		return "PDF"
	case SVG:
		return "SVG"
	case MAGICK:
		return "MAGICK"
	case HEIF:
		return "HEIF"
	case AVIF:
		return "AVIF"
	default:
		return "Unknown"
	}
}

func (t ImageType) String() string {
	return strconv.Itoa(int(t))
}

const (
	JPEG ImageType = iota + 1
	WEBP
	PNG
	TIFF
	GIF
	PDF
	SVG
	MAGICK
	HEIF
	AVIF
)

type IGravity interface {
	GetGravity() string
	fmt.Stringer
}

type Gravity int

func (g Gravity) GetGravity() string {
	switch g {
	case GravityCentre:
		return "Centre"
	case GravityNorth:
		return "North"
	case GravityEast:
		return "East"
	case GravitySouth:
		return "South"
	case GravityWest:
		return "West"
	case GravitySmart:
		return "Smart"
	default:
		return "Unknown"
	}
}

func (g Gravity) String() string {
	return strconv.Itoa(int(g))
}

const (
	GravityCentre Gravity = iota
	GravityNorth
	GravityEast
	GravitySouth
	GravityWest
	GravitySmart
)
