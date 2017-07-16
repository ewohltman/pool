package pool

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

// Set up a quick HTTP server for local testing
func init() {
	http.HandleFunc("/",
		func(respOut http.ResponseWriter, reqIn *http.Request) {
			defer reqIn.Body.Close()
		},
	)

	go func() {
		log.Fatalf("%s\n", http.ListenAndServe(":8080", nil))
	}()
}

// doPoolTest is the workhorse testing function
func doPoolTest(client Client) error {
	testURL, err := url.Parse("http://127.0.0.1:8080/")
	if err != nil {
		return err
	}

	var doFunc func(req *http.Request) (*http.Response, error)

	switch cp := client.(type) {
	case *PClient:
		doFunc = cp.DoPool
	case *http.Client:
		doFunc = cp.Do
	}

	resp, err := doFunc(&http.Request{URL: testURL})
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return nil
}

// TestNewPClient tests creating a PClient
func TestNewPClient(t *testing.T) {
	t.Logf("[INFO] Starting TestNewPClient")

	standardLibClient := &http.Client{}

	_ = NewPClient(standardLibClient, 25, 200) // Max 25 connections, 200 requests-per-second
	_ = NewPClient(standardLibClient, 0, 0)    // Why do this? Just use http.Client
	_ = NewPClient(standardLibClient, -1, -1)  // What

	t.Logf("[INFO] Completed TestNewPClient")
}

// TestPClient_Do tests performing a drop-in http.Client with pooling
func TestPClient_Do(t *testing.T) {
	t.Logf("[INFO] Starting TestPClient_Do")
	pClient := NewPClient(&http.Client{}, 0, 0)

	if err := doPoolTest(pClient); err != nil {
		t.Error("pool: ", err)
	}
	t.Logf("[INFO] Completed TestPClient_Do")
}

// TestPClient_DoPool tests performing a request with the pooling logic
func TestPClient_DoPool(t *testing.T) {
	t.Logf("[INFO] Starting TestPClient_DoPool")
	pClient := NewPClient(&http.Client{}, 0, 0)

	if err := doPoolTest(pClient); err != nil {
		t.Error("pool: ", err)
	}

	pClient2 := NewPClient(&http.Client{}, 25, 200)

	if err := doPoolTest(pClient2); err != nil {
		t.Error("pool: ", err)
	}

	pClient3 := NewPClient(&http.Client{}, -1, -1)

	if err := doPoolTest(pClient3); err != nil {
		t.Error("pool: ", err)
	}
	t.Logf("[INFO] Completed TestPClient_DoPool")
}

// BenchmarkPClient_Do benchmarks the pooling logic
func BenchmarkPClient_Do(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_Do")
	pClient := NewPClient(&http.Client{}, 0, 0)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_Do")
}

// BenchmarkPClient_DoPool benchmarks the pooling logic
func BenchmarkPClient_DoPool(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool")
	pClient := NewPClient(&http.Client{}, 0, 0)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool")
}

// BenchmarkBaseline benchmarks a request with the standard library http.Client
func BenchmarkBaseline(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkBaseline")
	stdClient := &http.Client{}

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(stdClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkBaseline")
}

// BenchmarkPClient_DoPool_10_0 benchmarks the pooling logic
func BenchmarkPClient_DoPool_10_0(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_10_0")
	pClient := NewPClient(&http.Client{}, 10, 0)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_10_0")
}

// BenchmarkPClient_DoPool_0_10 benchmarks the pooling logic
func BenchmarkPClient_DoPool_0_10(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_0_10")
	pClient := NewPClient(&http.Client{}, 0, 10)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_0_10")
}

// BenchmarkPClient_DoPool_10_10 benchmarks the pooling logic
func BenchmarkPClient_DoPool_10_10(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_10_10")
	pClient := NewPClient(&http.Client{}, 10, 10)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_10_10")
}

// BenchmarkPClient_DoPool_10_100 benchmarks the pooling logic
func BenchmarkPClient_DoPool_10_100(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_10_100")
	pClient := NewPClient(&http.Client{}, 10, 100)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_10_100")
}

// BenchmarkPClient_DoPool_10_200 benchmarks the pooling logic
func BenchmarkPClient_DoPool_10_200(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_10_200")
	pClient := NewPClient(&http.Client{}, 10, 200)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_10_200")
}

// BenchmarkPClient_DoPool_20_100 benchmarks the pooling logic
func BenchmarkPClient_DoPool_20_100(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_20_100")
	pClient := NewPClient(&http.Client{}, 20, 100)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_20_100")
}

// BenchmarkPClient_DoPool_20_200 benchmarks the pooling logic
func BenchmarkPClient_DoPool_20_200(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_20_200")
	pClient := NewPClient(&http.Client{}, 20, 200)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_20_200")
}

// BenchmarkPClient_DoPool_30_100 benchmarks the pooling logic
func BenchmarkPClient_DoPool_30_100(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_30_100")
	pClient := NewPClient(&http.Client{}, 30, 100)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_30_100")
}

// BenchmarkPClient_DoPool_30_200 benchmarks the pooling logic
func BenchmarkPClient_DoPool_30_200(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool_30_200")
	pClient := NewPClient(&http.Client{}, 30, 200)

	for n := 0; n < b.N; n++ {
		if err := doPoolTest(pClient); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool_30_200")
}
