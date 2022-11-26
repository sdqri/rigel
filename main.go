package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/sdqri/rigel/adapters"
	"github.com/sdqri/rigel/config"
	ctrl "github.com/sdqri/rigel/controller"
	"github.com/sdqri/rigel/service"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"

	_ "github.com/sdqri/rigel/docs"
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
		"appname":  "rigel",
		"pid":      strconv.Itoa(pid),
	})
	main_entry := entry.WithFields(log.Fields{
		"package": "main",
	})
	main_entry.Debug("Into this world, we're thrown!")

	config := config.GetConfig()

	redisAdp := adapters.NewRedisClient(
		main_entry,           //LogEntry
		config.Prefix,        //Prefix
		config.RedisAddress,  //addr
		config.RedisPassword, //password
		config.RedisDB,       //db
		5*time.Second,        //timeout
	)

	memAdp := adapters.NewMemoryClient[*service.RemoteImage](
		entry,         //logEntry
		config.Prefix, //Prefix
		10,            //cap
	)

	// Creating Algkey
	pub, err := jwk.ParseKey(config.PubKeyPem, jwk.WithPEM(true))
	if err != nil {
		logger.Fatalf("Error while trying create AlgKey, err = %v", err)
	}
	algKey := service.AlgKey{
		Alg:    jwa.SignatureAlgorithm(config.Alg),
		PubKey: pub,
	}

	controller := ctrl.New(
		entry,          //logEntry
		false,          //debug
		config.Prefix,  //Prefix
		config.Version, //Version
		algKey,         //AlgKey
		[]adapters.Cacher{memAdp, redisAdp},
	)

	server := fiber.New()
	server.Use(recover.New())

	server.Get("/docs/*", swagger.HandlerDefault)
	server.Get("/docs/*", swagger.New(swagger.Config{
		URL:          fmt.Sprintf("%s:%v/openapi.json", config.Host, config.Port),
		DocExpansion: "none",
	}))

	server.Mount(config.Prefix, controller.App)
	addr := fmt.Sprintf("%s:%v", config.Host, config.Port)
	server.Listen(addr)
}
