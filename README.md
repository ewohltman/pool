# pool - An HTTP Client with autonomous connection pooling and rate limiting
[![GoDoc](https://godoc.org/github.com/ewohltman/pool?status.svg)](https://godoc.org/github.com/ewohltman/pool)
[![Go Report Card](https://goreportcard.com/badge/github.com/ewohltman/pool)](https://goreportcard.com/report/github.com/ewohltman/pool)
[![Build Status](https://travis-ci.org/ewohltman/pool.svg?branch=master)](https://travis-ci.org/ewohltman/pool)

<br/>

`pool` wraps a standard `*http.Client` to add the ability to put a maximum number of connections in the pool for the client and the requests-per-second it can perform requests at.

It 'overloads' the `http.Client.Do(req *http.Request) (*http.Response, error)` method to implement the extended functionality.  By doing so, existing codebases do not need to heavily re-factor how they already do their logic to see the effect of the more finely-tuned client.

Functions that take in a `pool.Client` allow for the ability to take in either a `*pool.PClient` or an `*http.Client` (since they both implement `Do(req *http.Request) (*http.Response, error)`).  The function must operate on the argument's Do method since it will satisfy the interface for all types passed into it.  See `doPoolTest(client Client) error` in pool_test.go for an example.

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
