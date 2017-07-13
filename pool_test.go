// Tests for pool.go

package pool

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

// Set up a quick HTTP server for local testing
func init() {
	http.HandleFunc("/",
		func(respOut http.ResponseWriter, reqIn *http.Request) {
			defer reqIn.Body.Close()

			return
		},
	)

	go http.ListenAndServe(":8080", nil)
}

// TestNewPClient tests creating a PClient
func TestNewPClient(t *testing.T) {
	standardLibClient := &http.Client{}

	// pooledClient := NewPClient(standardLibClient, 25, 200) // Max 25 connections, 200 requests-per-second
	_ = NewPClient(standardLibClient, 25, 200)

	// normalClient := NewPClient(standardLibClient, 0, 0) // Why do this? Just use http.Client
	_ = NewPClient(standardLibClient, 0, 0)

	return
}

// TestPClient_Do tests performing a drop-in http.Client with pooling
func TestPClient_Do(t *testing.T) {
	if err := doTest(false); err != nil {
		t.Error("pool: ", err)
	}

	return
}

// TestPClient_DoPool tests performing a request with the pooling logic
func TestPClient_DoPool(t *testing.T) {
	if err := doTest(true); err != nil {
		t.Error("pool: ", err)
	}

	return
}

// doTest performs a standard GET request against the local HTTP server
func doTest(pool bool) error {
	standardLibClient := &http.Client{}

	pClient := NewPClient(standardLibClient, 25, 200)

	testURL, err := url.Parse("http://127.0.0.1/")
	if err != nil {
		return err
	}

	req := &http.Request{URL: testURL}

	var testFunc func(req *http.Request) (*http.Response, error)

	if pool {
		testFunc = pClient.DoPool
	} else {
		testFunc = pClient.Do
	}

	resp, err := testFunc(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
