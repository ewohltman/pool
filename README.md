# pool - An HTTP Client with autonomous connection pooling and rate limiting
[![GoDoc](https://godoc.org/github.com/ewohltman/pool?status.svg)](https://godoc.org/github.com/ewohltman/pool)
[![Go Report Card](https://goreportcard.com/badge/github.com/ewohltman/pool)](https://goreportcard.com/report/github.com/ewohltman/pool)

## Installation

- - -

`go get -u github.com/ewohltman/pool`

## Usage

- - -

```go
// Some other examples are in pool_test.go

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/ewohltman/pool"
)

func main() {
    standardLibClient := &http.Client{}
    
    pooledClient := pool.NewPClient(standardLibClient, 25, 200)
    
    urlString := "https://yourFavoriteWebsite.com/"
    
    reqURL, err := url.Parse(urlString)
    if err != nil {
        fmt.Printf("[ERROR] Unable to parse: %s", urlString)
        os.Exit(1)
    }
    
    resp, err := pooledClient.Do(&http.Request{URL: reqURL})
    if err != nil {
        fmt.Printf("[ERROR] Unable to perform request: %s", err)
        os.Exit(2)
    }
    defer resp.Body.Close()
    
    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("[ERROR] Unable to read response body: %s", err)
        os.Exit(3)
    }
}

```
