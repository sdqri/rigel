package gorigelsdk_test

import (
	"encoding/json"
	"fmt"
	"testing"

	rigelsdk "github.com/sdqri/rigel/sdk"
	"github.com/stretchr/testify/assert"
)

func TestProxyImageWithoutOptionsAndExpiry(t *testing.T) {
	key := "secretkey"
	salt := "secretsalt"
	rigelSDK := rigelsdk.NewSDK("http://localhost:8080/rigel", key, salt)

	imageURL := "https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg"

	actual := rigelSDK.ProxyImage(imageURL, nil, -1)

	expected := "http://localhost:8080/rigel/proxy?img=https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg&X-Signature=vX59TgdwdNqZD_jXGOky_zVgttc"
	assert.Equal(t, expected, actual)
}

func TestProxyImageWithOptionsWithoutExpiry(t *testing.T) {
	key := "secretkey"
	salt := "secretsalt"
	rigelSDK := rigelsdk.NewSDK("http://localhost:8080/rigel", key, salt)

	imageURL := "https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg"
	expectedSignature := "zkEmP1FDNoopC8GoM-caGzx1_1s"

	width, height := 100, 100
	options := rigelsdk.Options{
		Width:  &width,
		Height: &height,
		Type:   rigelsdk.WEBP,
	}

	actual := rigelSDK.ProxyImage(imageURL, &options, -1)

	expected := fmt.Sprintf("http://localhost:8080/rigel/proxy?height=100&img=%s&type=2&width=100&X-Signature=%s", imageURL, expectedSignature)
	assert.Equal(t, expected, actual)
}

func TestProxyImageWithOptionsAndExpiry(t *testing.T) {
	key := "secretkey"
	salt := "secretsalt"
	rigelSDK := rigelsdk.NewSDK("http://localhost:8080/rigel", key, salt)

	imageURL := "https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg"

	expectedSignature := "v6ROumbVPw18CkoBk9auEktWlzo"

	width, height := 100, 100
	options := rigelsdk.Options{
		Width:  &width,
		Height: &height,
		Type:   rigelsdk.WEBP,
	}

	actual := rigelSDK.ProxyImage(imageURL, &options, 1000*60*60*24)

	expected := fmt.Sprintf("http://localhost:8080/rigel/proxy?X-ExpiresAt=86400000&height=100&img=%s&type=2&width=100&X-Signature=%s", imageURL, expectedSignature)
	assert.Equal(t, expected, actual)
}

func TestCacheImage(t *testing.T) {
	key := "secretkey"
	salt := "secretsalt"
	rigelSDK := rigelsdk.NewSDK("http://localhost:8080/rigel", key, salt)

	imageURL := "https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg"
	expected := "http://localhost:8080/rigel/img/fde5eda7214568293ad70621aec2ad1efee5c7fd?X-Signature=ztW09e3EvM5IE7fJNsg0Z5-lPXg"

	options := rigelsdk.Options{
		Width:  rigelsdk.Ptr(300),
		Height: rigelsdk.Ptr(300),
		Type:   rigelsdk.WEBP,
	}

	actual, err := rigelSDK.CacheImage(imageURL, &options, -1)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestBatchedCacheImage(t *testing.T) {
	key := "secretkey"
	salt := "secretsalt"
	rigelSDK := rigelsdk.NewSDK("http://localhost:8080/rigel", key, salt)

	batchedCachedImageArgs := []rigelsdk.ProxyParams{
		{
			Img: "https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg",
			Options: rigelsdk.Options{
				Height: rigelsdk.Ptr(100),
				Width:  rigelsdk.Ptr(100),
				Type:   rigelsdk.WEBP,
			},
		},
		{
			Img: "https://img.freepik.com/premium-photo/baby-cat-british-shorthair_648604-47.jpg",
			Options: rigelsdk.Options{
				Height: rigelsdk.Ptr(100),
				Width:  rigelsdk.Ptr(100),
				Type:   rigelsdk.WEBP,
			},
		},
	}

	expected := `[{"img":"https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg","signature":"124799fa1f5d2069e1b56793e01f8fe260b87791"},{"img":"https://img.freepik.com/premium-photo/baby-cat-british-shorthair_648604-47.jpg","signature":"7fba571dee9007af7964e23239e2a1201419c0b8"}]`
	result, err := rigelSDK.BatchedCacheImage(batchedCachedImageArgs, -1)
	assert.Nil(t, err)

	actualBytes, err := json.Marshal(result)
	actual := string(actualBytes)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestTryShortURL(t *testing.T) {
	key := "secretkey"
	salt := "secretsalt"
	rigelSDK := rigelsdk.NewSDK("http://localhost:8080/rigel", key, salt)

	const imageURL = "https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg"
	const shortURLExpected = "http://localhost:8080/rigel/img/fde5eda7214568293ad70621aec2ad1efee5c7fd?X-Signature=ztW09e3EvM5IE7fJNsg0Z5-lPXg"

	width, height := 300, 300
	options := rigelsdk.Options{
		Width:  &width,
		Height: &height,
		Type:   rigelsdk.WEBP,
	}

	shortURL := rigelSDK.TryShortURL(imageURL, &options, -1)

	assert.Equal(t, shortURLExpected, shortURL, "Short URL mismatch")
}
