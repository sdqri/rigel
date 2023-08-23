<p align="center">
 <h2 align="center">rigelsdk</h2>
 <p align="center">Typescript SDK for <a href="[rigel](https://github.com/sdqri/rigel)">rigel</a></p>
  <p align="center">
  <a href="https://github.com/sdqri/gorigelsdk/issues">
      <img alt="Issues" src="https://img.shields.io/github/issues/sdqri/gorigelsdk?style=flat&color=336791" />
    </a>
    <a href="https://github.com/sdqri/gorigelsdk/pulls">
      <img alt="GitHub pull requests" src="https://img.shields.io/github/issues-pr/sdqri/gorigelsdk?style=flat&color=336791" />
    </a>
     <a href="https://github.com/sdqri/gorigelsdk">
      <img alt="GitHub Downloads" src="https://img.shields.io/npm/dw/gorigelsdk?style=flat&color=336791" />
    </a>
    <a href="https://github.com/sdqri/gorigelsdk">
      <img alt="GitHub Total Downloads" src="https://img.shields.io/npm/dt/gorigelsdk?color=336791&label=Total%20downloads" />
    </a>
 <a href="https://github.com/sdqri/gorigelsdk">
      <img alt="GitHub release" src="https://img.shields.io/github/release/sdqri/gorigelsdk.svg?style=flat&color=336791" />
    </a>
    <br />
  <a href="https://github.com/sdqri/gorigelsdk/issues/new/choose">Report Bug</a> /
  <a href="https://github.com/sdqri/gorigelsdk/issues/new/choose">Request Feature</a>
  </p>

# Getting started

## Installation

> Install using Go modules:

```bash
go get github.com/sdqri/gorigelsdk
```

### Import the sdk

```
import "github.com/sdqri/gorigelsdk"
```

### Usage examples

```go
package main

import (
	"fmt"
	"time"

	gorigelsdk "github.com/sdqri/gorigelsdk"
)

func main() {
	KEY := "secretkey"
	SALT := "secretsalt"
	BASE_URL := "<put rigel url here>" // e.g., http://localhost:8080/rigel

	rigelSDK := gorigelsdk.NewSDK(BASE_URL, KEY, SALT)

	proxyURL, _ := rigelSDK.ProxyImage(
		"https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg",
		&gorigelsdk.Options{Width: 100, Height: 100, Type: gorigelsdk.ImageTypeWEBP},
		time.Now().Add(24*time.Hour).Unix(), // 1 day expiry
	)

	// Creating short URL
	shortURL, _ := rigelSDK.TryShortURL(
		"https://www.pakainfo.com/wp-content/uploads/2021/09/image-url-for-testing.jpg",
		&gorigelsdk.Options{Width: 100, Height: 100, Type: gorigelsdk.ImageTypeWEBP},
		time.Now().Add(24*time.Hour).Unix(), // 1 day expiry
	)

	fmt.Println(proxyURL)
	fmt.Println(shortURL)
}

```

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](issues).

## Show your support

Give a ⭐️ if this project helped you!

## 📝 License

This SDK is distributed under the
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0),
see LICENSE.txt for more information.
