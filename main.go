package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sdqri/rigel/adapters"
	"github.com/sdqri/rigel/config"
	ctrl "github.com/sdqri/rigel/controllers"
	"github.com/sdqri/rigel/middlewares"
	srv "github.com/sdqri/rigel/services"
	"github.com/sdqri/rigel/utils"

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
		signatory := utils.NewSignatory(config.Prefix, config.XKey, config.XSalt)
		signatureValidator = middlewares.NewSignatureValidator(config.Prefix, signatory)
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
	fmt.Printf("Listening and serving HTTP on %v", server_addr)
	srv := &http.Server{
		Addr:    server_addr,
		Handler: router,
	}

	go func() {
		// service connections
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("unable to ListenAndServe, err: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
