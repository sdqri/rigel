package controller

import (
	"errors"
	"fmt"

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
	Prefix   string
	Version  string
	*fiber.App
	cachers []adapters.Cacher
}

func New(logEntry *log.Entry, prefix, version string, cashers []adapters.Cacher, fiberConfig ...fiber.Config) *RigelController {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.controller",
	})

	router := fiber.New(fiberConfig...)
	ctrl := &RigelController{
		LogEntry: entry,
		Prefix:   prefix,
		Version:  version,
		App:      router,
		cachers:  cashers,
	}

	ctrl.Get("/version", ctrl.getVersion)
	ctrl.Get("/", ctrl.getImage)
	return ctrl
}

func (ctrl *RigelController) getVersion(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"version": ctrl.Version})
}

func (ctrl *RigelController) getImage(c *fiber.Ctx) error {
	// Checking whether query parameter exists
	queryParams := c.Query("req")
	if queryParams == "" {
		return ErrNoQueryParameters
	}

	// Parsing RemoteImage for finding src
	imageRequest, err := service.ParseToken(queryParams)
	if err != nil {
		return err
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
	remoteImage, err = imageRequest.Download()
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
	c.Attachment(fileName)
	return c.Send(*remoteImage.Data)
}
