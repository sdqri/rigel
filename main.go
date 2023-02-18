package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sdqri/rigel/adapters"
	"github.com/sdqri/rigel/config"
	ctrl "github.com/sdqri/rigel/controllers"
	"github.com/sdqri/rigel/middlewares"
	srv "github.com/sdqri/rigel/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Creating Root Log.Entry ------------------------------------------------
	logger := log.New()
	logger.SetLevel(log.DebugLevel) // DebugLevel for verbose logging
	logger.SetFormatter(&log.JSONFormatter{})
	hostname, err := os.Hostname()
	if err != nil {
		logger.Debugf("Error while trying to get host name, err = %v", err)
		hostname = "error"
	}
	pid := os.Getpid()
	entry := logger.WithFields(log.Fields{
		"hostname": hostname,
		"appname":  "rigel",
		"pid":      strconv.Itoa(pid),
	})
	main_entry := entry.WithFields(log.Fields{
		"package": "main",
	})
	main_entry.Debug("Into this world, we're thrown!")

	config := config.GetConfig()

	redisAdp := adapters.NewRedisClient(
		main_entry,           // LogEntry
		config.Prefix,        // Prefix
		config.RedisAddress,  // addr
		config.RedisPassword, // password
		config.RedisDB,       // db
		time.Duration(config.RedisTimeout)*time.Second,    //timeout
		time.Duration(config.RedisExpiration)*time.Second, //expiration
	)

	memAdp := adapters.NewMemoryClient(
		entry,         // logEntry
		config.Prefix, // p[*srv.RemoteImage]refix
		config.Cap,    // cap
	)

	rigelService := srv.NewRigelService(config.Debug, entry,
		adapters.NewMultilevelCacher[string](memAdp, redisAdp))

	var signatureValidator gin.HandlerFunc = nil
	if config.SignatureValidation == true {
		signatureValidator = middlewares.NewSignatureValidator(config.XKey, config.XSalt, config.Prefix)
	}

	rigelController := ctrl.NewRigelController(
		config.Debug,       // debug
		entry,              // logEntry
		config.Version,     // version
		rigelService,       // service
		signatureValidator, // signatureValidator
	)

	router := gin.Default()
	if config.CORS == true {
		router.Use(cors.New(cors.Config{
			AllowOrigins: config.AllowOrigins,
			AllowMethods: config.AllowMethods,
			AllowHeaders: config.AllowHeaders,
		}))
	}

	prefix := router.Group(config.Prefix)
	rigelController.Handle(prefix)

	server_addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	router.Run(server_addr)
}
