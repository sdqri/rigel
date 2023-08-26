package specifications

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type E2EAPI interface {
	Version() Result[string]
	Proxy(pathParams map[string]string) Result[*[]byte]
}

var ENVS = map[string]string{
	"PUBLIC_KEY_PEM": "test_key",
}

func setupServer(logger *zap.Logger) func() {
	logger = logger.With(zap.String("function", "setupServer"))
	logger.Debug("Setting up test suit")

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal("unable to get working directory", zap.Error(err))
	}

	// Run setup_suit
	cmd := exec.Command("./setup_suit.sh")
	cmd.Dir = path.Join(wd, "test_fixtures")
	_, err = cmd.Output()
	if err != nil {
		logger.Fatal("failed running setup_suit.sh", zap.Error(err))
	}
	logger.Debug("setup_suit run is successfull")

	// Build project
	buildCmd := exec.Command("go", "build", ".")
	buildCmd.Dir = filepath.Dir(wd)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Run()
	if err != nil {
		logger.Fatal("unable to build main.go", zap.Error(err))
	}

	// Run rigel
	serverCtx, cancelServer := context.WithCancel(context.Background())
	serverCmd := exec.CommandContext(serverCtx, "./rigel")
	serverCmd.Dir = filepath.Dir(wd) //parent of specifications directory
	serverCmd.Env = os.Environ()
	for envVar, envVal := range ENVS {
		serverCmd.Env = append(serverCmd.Env, fmt.Sprintf("%s=%s", envVar, envVal))
	}

	outputBuffer := bytes.NewBuffer(nil)
	serverCmd.Stdout = outputBuffer

	successChan := make(chan struct{}, 0)
	go func(buf *bytes.Buffer, successChan chan struct{}) {
		for {
			if strings.Contains(buf.String(), "Listening and serving HTTP") {
				successChan <- struct{}{}
				break
			}
		}
	}(outputBuffer, successChan)

	err = serverCmd.Start()
	if err != nil {
		logger.Fatal("failed running main.go", zap.Error(err))
	}
	<-successChan
	logger.Debug(":)")

	return func() {
		logger.Debug("Tearing down test suit")
		cancelServer()
		err := serverCmd.Process.Signal(syscall.SIGINT)
		cmd := exec.Command("./teardown_suit.sh")
		cmd.Dir = path.Join(wd, "test_fixtures")
		_, err = cmd.Output()
		if err != nil {
			logger.Fatal("failed running teardown_suit.sh", zap.Error(err))
		}
	}
}

var driver E2EAPI

// Define a custom flag to enable golden file update
var updateGolden = flag.Bool("update", false, "update golden files")

func TestMain(m *testing.M) {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = logger.With(zap.String("package", "TestMain"))

	if *updateGolden {
		logger = logger.With(zap.String("mode", "updating golden files"))
	}

	logger.Debug("Running test main")

	httpClient := http.Client{
		Timeout: 1 * time.Second,
	}
	driver = &E2EDriver{BaseURL: "http://localhost:8080/rigel", Client: &httpClient}

	teardownServer := setupServer(logger)

	exitCode := m.Run()
	teardownServer()

	os.Exit(exitCode)
}

func TestVersion(t *testing.T) {
	result := driver.Version()
	assert.NoError(t, result.Err)
	assert.NotEqual(t, result.Ok, "")
}

func TestProxy(t *testing.T) {

	pathParams := map[string]string{
		"img":         "https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg",
		"height":      "100",
		"width":       "100",
		"type":        "2",
		"X-Signature": "",
	}
	result := driver.Proxy(pathParams)
	assert.NoError(t, result.Err)
	actual := result.Ok

	if *updateGolden {
		// Update the golden file with the actual output
		err := os.WriteFile("test_fixtures/golden_test_proxy.webp", *actual, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Read the contents of the golden file
	expected, err := os.ReadFile("test_fixtures/golden_test_proxy.webp")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, *actual, expected)
}

// func Test(t *testing.T) {
// 	cases := []struct {
// 		Name             string
// 		UserCreateInList []models.UserCreateIn
// 		Expected         []models.UserOut
// 		Err              error
// 	}{
// 		{
// 			"list_of_length2",
// 			[]models.UserCreateIn{
// 				{
// 					FirstName: "fake_first_name0",
// 					LastName:  "fake_last_name0",
// 					Email:     "fake_email0@email.com",
// 					Password:  "F@ke_password0",
// 				},
// 				{
// 					FirstName: "fake_first_name1",
// 					LastName:  "fake_last_name1",
// 					Email:     "fake_email1@email.com",
// 					Password:  "F@ke_password1",
// 				},
// 			},
// 			[]models.UserOut{
// 				models.UserOut{
// 					FirstName: "fake_first_name0",
// 					LastName:  "fake_last_name0",
// 					Email:     "fake_email0@email.com",
// 					Password:  "F@ke_password001",
// 				},
// 				models.UserOut{
// 					FirstName: "fake_first_name1",
// 					LastName:  "fake_last_name1",
// 					Email:     "fake_email1@email.com",
// 					Password:  "F@ke_password001",
// 				},
// 			},
// 			nil,
// 		},
// 	}
//
// 	for _, tc := range cases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			UserOutList := make([]models.UserOut, 0)
// 			for _, userCreatedIn := range tc.UserCreateInList {
// 				res := UserClient.Create(userCreatedIn)
// 				assert.NoError(t, res.Err)
// 				UserOutList = append(UserOutList, res.Ok)
// 			}
// 			// Only testing length of Expected List
// 			assert.Equal(t, len(tc.Expected), len(UserOutList))
// 		})
// 	}
// }
