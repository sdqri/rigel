package service

type ImageService interface {
	Process() ([]byte, error)
}