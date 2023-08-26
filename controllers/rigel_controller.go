package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sdqri/rigel/services"

	log "github.com/sirupsen/logrus"
)

var (
	ErrNoQueryParameters error = errors.New("no req query string parameter in url")
)

type ProxyParams struct {
	Img       string             `json:"img" form:"img" binding:"required"`
	IROptions services.IROptions `json:"options"`
}

type CacheImageResponse struct {
	Img       string `json:"img"`
	Signature string `json:"signature"`
}

type ImageService interface {
	ProxyImageRequest(imageRequest *services.ImageRequest) (*services.RemoteImage, error)
	CacheImageRequest(imageRequest *services.ImageRequest) error
	GetBySignature(string) (*services.RemoteImage, error)
}

type RigelController struct {
	debug              bool
	LogEntry           *log.Entry
	Version            string
	Service            ImageService
	SignatureValidator gin.HandlerFunc
}

func NewRigelController(
	debug bool, logEntry *log.Entry,
	version string, service ImageService,
	signatureValidator gin.HandlerFunc) *RigelController {
	// Setting package specific fields for log entry
	entry := logEntry.WithFields(log.Fields{
		"package": "adapters.rigel_controller",
	})

	ctrl := &RigelController{
		debug:              debug,
		LogEntry:           entry,
		Version:            version,
		Service:            service,
		SignatureValidator: signatureValidator,
	}
	return ctrl
}

func (ctrl *RigelController) Handle(router gin.IRouter) {

	if ctrl.SignatureValidator != nil {
		router.GET("/version", ctrl.getVersion)
		router.GET("/proxy", ctrl.SignatureValidator, ctrl.ProxyImage)
		router.POST("/headsup", ctrl.SignatureValidator, ctrl.CacheImage)
		router.POST("/batched-headsup", ctrl.SignatureValidator, ctrl.BatchedCacheImage)
		router.GET("/img/:signature", ctrl.SignatureValidator, ctrl.GetBySignature)
	} else {
		router.GET("/version", ctrl.getVersion)
		router.GET("/proxy", ctrl.ProxyImage)
		router.POST("/headsup", ctrl.CacheImage)
		router.POST("/batched-headsup", ctrl.BatchedCacheImage)
		router.GET("/img/:signature", ctrl.GetBySignature)
	}
}

func (ctrl *RigelController) getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"version": ctrl.Version})
}

func (rc *RigelController) ProxyImage(c *gin.Context) {
	var args ProxyParams
	err := c.Bind(&args)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	imageRequest, err := services.NewImageRequest(args.Img, args.IROptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	remoteImage, err := rc.Service.ProxyImageRequest(imageRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.Data(http.StatusOK, remoteImage.ContentType(), *remoteImage.Data)
}

func (rc *RigelController) CacheImage(c *gin.Context) {
	var args ProxyParams
	err := c.BindQuery(&args)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	imageRequest, err := services.NewImageRequest(args.Img, args.IROptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = rc.Service.CacheImageRequest(imageRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	signature, err := imageRequest.SHA1Sum()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Fire ProxyImageRequest to download, process and cache RemoteImage
	go rc.Service.ProxyImageRequest(imageRequest)

	c.JSON(http.StatusOK, CacheImageResponse{Img: args.Img, Signature: signature})
}

func (rc *RigelController) GetBySignature(c *gin.Context) {
	signature := c.Param("signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "signature hasn't provided",
		})
		return
	}

	remoteImage, err := rc.Service.GetBySignature(signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "signature isn't valid",
		})
		return
	}

	c.Data(http.StatusOK, remoteImage.ContentType(), *remoteImage.Data)
}

func (rc *RigelController) BatchedCacheImage(c *gin.Context) {
	var args []ProxyParams
	if err := c.Bind(&args); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	response := make([]CacheImageResponse, 0)
	for _, proxyParams := range args {
		imageRequest, err := services.NewImageRequest(proxyParams.Img, proxyParams.IROptions)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		err = rc.Service.CacheImageRequest(imageRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		go rc.Service.ProxyImageRequest(imageRequest)

		signature, err := imageRequest.SHA1Sum()
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		response = append(response, CacheImageResponse{Img: proxyParams.Img, Signature: signature})
	}

	c.JSON(http.StatusOK, response)
}
