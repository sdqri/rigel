package controller

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/sdqri/rigel/adapters"
	"github.com/sdqri/rigel/service"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNoQueryParameters error = errors.New("no req query string parameter in url")
)

type RigelController struct {
	LogEntry *log.Entry
	debug    bool
	Prefix   string
	Version  string
	service.AlgKey
	*fiber.App
	cachers     []adapters.Cacher
	redisClient *adapters.RedisClient
}

func New(logEntry *log.Entry, debug bool, prefix, version string, algKey service.AlgKey, cashers []adapters.Cacher, redisClient *adapters.RedisClient, fiberConfig ...fiber.Config) *RigelController {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.controller",
	})

	router := fiber.New(fiberConfig...)
	ctrl := &RigelController{
		LogEntry:    entry,
		debug:       debug,
		Prefix:      prefix,
		Version:     version,
		AlgKey:      algKey,
		App:         router,
		cachers:     cashers,
		redisClient: redisClient,
	}

	ctrl.Get("/version", ctrl.getVersion)
	ctrl.Get("/img/:req", ctrl.getImage)
	ctrl.Get("/headsup/:req", ctrl.headImage)
	return ctrl
}

func (ctrl *RigelController) getVersion(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"version": ctrl.Version})
}

func (ctrl *RigelController) getImage(c *fiber.Ctx) error {
	// Checking whether query parameter exists
	queryParams := c.Params("req", "noreq")
	if queryParams == "noreq" {
		c.SendStatus(404)
		return c.SendString(ErrNoQueryParameters.Error())
	}

	var imageRequest *service.ImageRequest
	// SHA-1 address | short URL
	if len(queryParams) <= 64 {
		imageRequest = service.NewImageRequest(queryParams)
	} else {
		// Parsing RemoteImage for finding src
		var err error
		imageRequest, err = service.ParseToken(ctrl.AlgKey, queryParams, ctrl.debug)
		if err != nil {
			c.SendStatus(404)
			return c.SendString(err.Error())
		}
	}

	var remoteImage *service.RemoteImage
	remoteImage = service.NewRemoteImage(service.WithImageRequest(imageRequest))
	// check chachers
	for _, cacher := range ctrl.cachers {
		err := cacher.GetCachable(remoteImage)
		if err == nil {
			fileName := fmt.Sprintf("image.%s", remoteImage.Type())
			c.Attachment(fileName)
			return c.Send(*remoteImage.Data)
		}
	}

	// Downloading image
	remoteImage, err := imageRequest.Download()
	if err != nil {
		return fiber.NewError(fiber.StatusFailedDependency, "error while trying to download image")
	}

	// Processing image
	err = remoteImage.Process()
	if err != nil {
		return fiber.NewError(fiber.StatusFailedDependency, "error while processing image")
	}

	// Caching
	go func() {
		for _, cacher := range ctrl.cachers {
			cacher.Cache(remoteImage)
		}
	}()

	fileName := fmt.Sprintf("image.%s", remoteImage.Type())
	c.Type(filepath.Ext(fileName))
	// TODO: Add inline content disposition
	// c.setCanonical(HeaderContentDisposition, `attachment; filename="`+c.app.quoteString(fname)+`"`)
	// c.Attachment(fileName)
	return c.Send(*remoteImage.Data)
}

func (ctrl *RigelController) headImage(c *fiber.Ctx) error {
	// Checking whether query parameter exists
	queryParams := c.Params("req", "noreq")
	if queryParams == "noreq" {
		return c.SendStatus(404)
	}
	go func() {
		var imageRequest *service.ImageRequest
		// SHA-1 address | short URL
		if len(queryParams) <= 64 {
			imageRequest = service.NewImageRequest(queryParams)
		} else {
			// Parsing RemoteImage for finding src
			var err error
			imageRequest, err = service.ParseToken(ctrl.AlgKey, queryParams, ctrl.debug)
			if err != nil {
				return
			}
		}

		var remoteImage *service.RemoteImage
		remoteImage = service.NewRemoteImage(service.WithImageRequest(imageRequest))
		// check chachers
		for _, cacher := range ctrl.cachers {
			err := cacher.GetCachable(remoteImage)
			if err == nil {
				return
			}
		}

		// Downloading image
		remoteImage, err := imageRequest.Download()
		if err != nil {
			return
		}

		// Processing image
		err = remoteImage.Process()
		if err != nil {
			return
		}

		// Caching
		go func() {
			for _, cacher := range ctrl.cachers {
				cacher.Cache(remoteImage)
			}
		}()

		return
	}()
	return c.SendStatus(200)
}
