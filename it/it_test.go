package it_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var ENVS = map[string]string{
	"PUBLIC_KEY_PEM": "test_key",
}

func setupServer(logger *log.Logger) func() {
	// Logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic
	// Add context information to logger
	entry := logger.WithFields(log.Fields{
		"package": "setupServer",
	})
	entry.Debugf("Setting up test suit")

	wd, err := os.Getwd()
	if err != nil {
		entry.Fatalf("err : %v", err)
	}

	cmd := exec.Command("./setup_suit.sh")
	cmd.Dir = path.Join(wd, "test_fixtures")
	stdoe, err := cmd.Output()
	if err != nil {
		entry.Fatalf("failed running setup_suit.sh : %v \n %v", err, string(stdoe))
	}
	entry.Debugf("%v", string(stdoe))

	// setup server
	serverCtx, cancelServer := context.WithCancel(context.Background())
	serverCmd := exec.CommandContext(serverCtx, "go", "run", "main.go")
	serverCmd.Dir = filepath.Dir(wd) //parent of it directory
	serverCmd.Env = os.Environ()
	for envVar, envVal := range ENVS {
		serverCmd.Env = append(serverCmd.Env, fmt.Sprintf("%s=%s", envVar, envVal))
	}
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	err = serverCmd.Start()
	if err != nil {
		entry.Fatalf("failed running main.go : %v", err)
	}
	time.Sleep(30 * time.Second) //?

	return func() {
		entry.Debugf("Tearing down test suit")
		cancelServer()
		cmd := exec.Command("./teardown_suit.sh")
		cmd.Dir = path.Join(wd, "test_fixtures")
		stdoe, err := cmd.Output()
		if err != nil {
			entry.Fatalf("failed running teardown_suit.sh : %v \n %v", err, string(stdoe))
		}
		entry.Debugf("%v", string(stdoe))
	}
}

func TestMain(m *testing.M) {
	// Creating root Log.Entry
	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{})

	logger.Debug("Running test main")
	teardownServer := setupServer(logger)
	defer teardownServer()

	m.Run()
}

func TestA(t *testing.T) {
	assert.Equal(t, 1, 2)
}
