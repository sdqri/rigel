package controller

import (
	"errors"
	"fmt"
	"rigel/adapters"
	"rigel/config"
	"rigel/service"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrNoQueryParameters error = errors.New("no req query string parameter in url")
)

type RigelController struct {
	*fiber.App
	cacher adapters.Cacher 
	cfg *config.Config
}


func New(cacher adapters.Cacher, cfg *config.Config, fiberConfig  ...fiber.Config) *RigelController{
	router := fiber.New(fiberConfig...)
	ctrl := &RigelController{
		App: router,
		cacher: cacher,
		cfg: cfg,
	}

	ctrl.Get("/version", ctrl.getVersion)
	ctrl.Get("/", ctrl.getImage)
	return ctrl
}

func (ctrl *RigelController) getVersion(c *fiber.Ctx) error{
	return c.JSON(map[string]string{"version": ctrl.cfg.Version})
}


func (ctrl *RigelController) getImage(c *fiber.Ctx) error{
	queryParams := c.Query("req")
	if queryParams == "" {
		return ErrNoQueryParameters
	} 

	remoteImage := service.NewRemoteImage(service.WithPrefix(ctrl.cfg.Prefix))
	err := remoteImage.ParseQueryParams(queryParams)
	if err!=nil{
		return err
	}


	// check if redisnil is error
	err = ctrl.cacher.GetCachable(remoteImage)
	if err != nil{
		remoteImage.Process()
	} 
	
	format, err := remoteImage.Type()
	if err!=nil{
		return err
	}
	filename := fmt.Sprintf("image.%s", format)

	c.Attachment(filename)

	// Caching
	go func ()  {
		ctrl.cacher.Cache(remoteImage)
	}()

	return c.Send(*remoteImage.Data)
}
