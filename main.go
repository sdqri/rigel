package main

import (
	"fmt"
	"os"
	"rigel/adapters"
	"rigel/config"
	ctrl "rigel/controller"
	"rigel/service"
	"strconv"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	_ "rigel/docs"
)

// @title Rigel Api
// @version 1.0.0
// @description Yet another image proxy
// @termsOfService http://swagger.io/terms/
// @contact.name Sadiq Rahmati
// @contact.email sadeg.r1@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
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
		"appname":  "argos",
		"pid":      strconv.Itoa(pid),
	})
	main_entry := entry.WithFields(log.Fields{
		"package": "main",
	})
	main_entry.Debug("Into this world, we're thrown!")
	
	config := config.GetConfig()
	rc := adapters.NewRedisClient(
		main_entry, config.RedisAddress,
		config.RedisPassword, config.RedisDB,
		5 * time.Second)
	memoryAdapter := adapters.NewMemoryClient[*service.RemoteImage](entry, rc, 10)
    controller := ctrl.New(memoryAdapter, config)
	
	server := fiber.New()

	server.Get("/docs/*", swagger.HandlerDefault)

	server.Get("/docs/*", swagger.New(swagger.Config{
		URL: fmt.Sprintf("%s:%v/openapi.json", config.Host, config.Port),
		DocExpansion: "none",
	}))

	server.Mount(config.Prefix, controller.App)
	addr := fmt.Sprintf("%s:%v", config.Host, config.Port)
    server.Listen(addr)
}