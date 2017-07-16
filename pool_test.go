// Tests for pool.go
package pool

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

// matrix is our test matrix of inputs.  matrix is initialized in the init function
var matrix *testMatrix

type testMatrix struct {
	Tests []*testInputs
}

type testInputs struct {
	maxPoolSize int
	reqPerSec   int
}

// Set up a quick HTTP server for local testing
func init() {
	matrix = &testMatrix{
		Tests: make([]*testInputs, 0),
	}

	// Add test inputs to the end of this
	matrix.Tests = append(matrix.Tests, &testInputs{0, 0})
	matrix.Tests = append(matrix.Tests, &testInputs{10, 0})
	matrix.Tests = append(matrix.Tests, &testInputs{0, 100})
	matrix.Tests = append(matrix.Tests, &testInputs{20, 200})
	matrix.Tests = append(matrix.Tests, &testInputs{30, 300})

	http.HandleFunc("/",
		func(respOut http.ResponseWriter, reqIn *http.Request) {
			defer reqIn.Body.Close()
		},
	)

	go func() {
		log.Fatalf("%s\n", http.ListenAndServe(":8080", nil))
	}()
}

// TestNewPClient tests creating a PClient
func TestNewPClient(t *testing.T) {
	t.Logf("[INFO] Starting TestNewPClient")

	standardLibClient := &http.Client{}

	// pooledClient := NewPClient(standardLibClient, 25, 200) // Max 25 connections, 200 requests-per-second
	_ = NewPClient(standardLibClient, 25, 200)

	// normalClient := NewPClient(standardLibClient, 0, 0) // Why do this? Just use http.Client
	_ = NewPClient(standardLibClient, 0, 0)

	t.Logf("[INFO] Completed TestNewPClient")
}

// BenchmarkBaseline benchmarks a request with the standard library http.Client
func BenchmarkBaseline(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkBaseline")
	for n := 0; n < b.N; n++ {
		if err := doBaselineTest(); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkBaseline")
}

// TestPClient_Do tests performing a drop-in http.Client with pooling
func TestPClient_Do(t *testing.T) {
	t.Logf("[INFO] Starting TestPClient_Do")
	if err := doPoolTest(t, false); err != nil {
		t.Error("pool: ", err)
	}
	t.Logf("[INFO] Completed TestPClient_Do")
}

// BenchmarkPClient_Do benchmarks the pooling logic
func BenchmarkPClient_Do(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_Do")
	for n := 0; n < b.N; n++ {
		if err := doPoolTest(b, false); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_Do")
}

// TestPClient_DoPool tests performing a request with the pooling logic
func TestPClient_DoPool(t *testing.T) {
	t.Logf("[INFO] Starting TestPClient_DoPool")
	if err := doPoolTest(t, true); err != nil {
		t.Error("pool: ", err)
	}
	t.Logf("[INFO] Completed TestPClient_DoPool")
}

// BenchmarkPClient_DoPool benchmarks the pooling logic
func BenchmarkPClient_DoPool(b *testing.B) {
	b.Logf("[INFO] Starting BenchmarkPClient_DoPool")
	for n := 0; n < b.N; n++ {
		if err := doPoolTest(b, true); err != nil {
			b.Error("pool: ", err)
		}
	}
	b.Logf("[INFO] Completed BenchmarkPClient_DoPool")
}

func doBaselineTest() error {
	standardLibClient := &http.Client{}

	testURL, err := url.Parse("http://127.0.0.1:8080/")
	if err != nil {
		return err
	}

	req := &http.Request{URL: testURL}

	resp, err := standardLibClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)

	return err
}

func doPoolTest(tb testing.TB, pool bool) error {
	if tester, ok := tb.(*testing.T); ok {
		tb = tester
	}

	if tester, ok := tb.(*testing.B); ok {
		tb = tester
	}

	tb.Logf("[INFO] doPoolTest")
	defer func() {
		tb.Logf("[INFO] Completed doPoolTest")
	}()

	for _, testRun := range matrix.Tests {
		thisRun := testRun
		standardLibClient := &http.Client{}

		tb.Logf("[INFO] doPoolTest maxPoolSize: %d, reqPerSec: %d", thisRun.maxPoolSize, thisRun.reqPerSec)

		pClient := NewPClient(standardLibClient, thisRun.maxPoolSize, thisRun.reqPerSec)

		testURL, err := url.Parse("http://127.0.0.1:8080/")
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

		_, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	return nil
}
