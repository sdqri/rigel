package services

import (
	"github.com/sdqri/rigel/adapters"

	log "github.com/sirupsen/logrus"
)

type RigelService struct {
	debug    bool
	LogEntry *log.Entry
	MLC      *adapters.MultilevelCacher[string]
}

func NewRigelService(debug bool, logEntry *log.Entry, mlc *adapters.MultilevelCacher[string]) *RigelService {
	return &RigelService{
		debug:    debug,
		LogEntry: logEntry,
		MLC:      mlc,
	}
}

func (rs *RigelService) ProxyImageRequest(imageRequest *ImageRequest) (remoteImage *RemoteImage, err error) {
	queryableRemoteImage, err := imageRequest.GetQueryableRemoteImage()
	if err != nil {
		return
	}

	err = rs.MLC.GetCachable(queryableRemoteImage)
	if err != nil {
		//if queryableRemoteImage doesn't exists in cache
		remoteImage, err = imageRequest.GetRemoteImage()
		rs.MLC.Cache(imageRequest)
		rs.MLC.Cache(remoteImage)
		return
	}
	remoteImage = queryableRemoteImage
	return
}

func (rs *RigelService) CacheImageRequest(imageRequest *ImageRequest) error {
	return rs.MLC.Cache(imageRequest)
}

func (rs *RigelService) GetBySignature(signature string) (remoteImage *RemoteImage, err error) {
	queryableImageRequest := NewQueryableImageRequest(signature)
	queryableRemoteImage, err := queryableImageRequest.GetQueryableRemoteImage()
	if err != nil {
		return
	}
	err = rs.MLC.GetCachable(queryableRemoteImage)
	if err != nil {
		//if queryableRemoteImage doesn't exists in cache
		err = rs.MLC.GetCachable(queryableImageRequest)
		if err != nil {
			return
		}
		return rs.ProxyImageRequest(queryableImageRequest)
	}
	remoteImage = queryableRemoteImage
	return
}
