// Package pool is a wrapper for Pooled HTTP clients.
// See net/http/client.go.
//
// This is the higher-level Pooled Client interface.
// The lower-level Client implementation is in client.go.
// The lowest-level implementation is in transport.go.
package pool

import (
	"net/http"
	"time"
)

// Client is a an interface common with http.Client
//
// Use it as a function argument when you want to take in an http.Client
// or a pool.PClient because you can operate on either type's Do method
// interchangeably.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// A PClient is a pooled HTTP client. A pre-configured *http.Client should
// be passed in to this package's NewPClient() function to get a PClient.
//
// PClients are http.Clients wrapped to extend higher-level functionality.
// A PClient can, and should, be used like a usual http.Client.
//
// Pooling functionality is handled autonomously under the hood.
// PClient can have an upper bound set for the number of connections it
// will make.  They can also be set to rate limit the requests-per-second
// the PClient can perform.
//
// PClients have a Do() method that may help make converting existing
// projects to use this easier. The implementation may not have to change
// to gain the benefits, only the initial type.
//
// The PClient's Transport typically has internal state (cached TCP
// connections), so PClients should be reused instead of created as
// needed. PClients are safe for concurrent use by multiple goroutines.
//
// For all the interesting details, see client.go.
type PClient struct {
	client       *http.Client
	maxPoolSize  int
	cSemaphore   chan int
	reqPerSecond int
	rateLimiter  <-chan time.Time
}

// NewPClient returns a *PClient that wraps an *http.Client and sets the
// maximum pool size as well as the requests per second as integers.
//
// A zero for maxPoolSize will set no limit
// A zero for reqPerSec will set no limit
func NewPClient(stdClient *http.Client, maxPoolSize int, reqPerSec int) *PClient {
	var semaphore chan int = nil
	if maxPoolSize > 0 {
		semaphore = make(chan int, maxPoolSize) // Buffered channel to act as a semaphore
	}

	var emitter <-chan time.Time = nil
	if reqPerSec > 0 {
		emitter = time.NewTicker(time.Second / time.Duration(reqPerSec)).C // x req/s == 1s/x req (inverse)
	}

	return &PClient{
		client:       stdClient,
		maxPoolSize:  maxPoolSize,
		cSemaphore:   semaphore,
		reqPerSecond: reqPerSec,
		rateLimiter:  emitter,
	}
}

// Do is a method that 'overloads' the http.Client Do method for drop-in
// compatibility with existing codebases.
func (c *PClient) Do(req *http.Request) (*http.Response, error) {
	return c.DoPool(req)
}

// DoPool does the required synchronization with channels to
// perform the request in accordance to the PClient's maximum pool size
// and requests-per-second rate limiter.
//
// Caller should close resp.Body when done reading from it.
//
// It is an exported function in case a direct call to this method is
// desired
func (c *PClient) DoPool(req *http.Request) (*http.Response, error) {
	if c.maxPoolSize > 0 {
		c.cSemaphore <- 1 // Grab a connection from our pool
		defer func() {
			<-c.cSemaphore // Defer release our connection back to the pool
		}()
	}

	if c.reqPerSecond > 0 {
		<-c.rateLimiter // Block until a signal is emitted from the rateLimiter
	}

	// Perform the normal request using the underlying http.Client
	resp, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	// resp.Body intentionally not closed

	return resp, nil
}
