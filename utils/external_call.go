package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type stop struct {
	error
}

func CallExternalAPI(req *http.Request) ([]byte, int, error) {
	timeOutSec, _ := strconv.Atoi(os.Getenv("TIMEOUTSEC"))
	attempts, _ := strconv.Atoi(os.Getenv("ATTEMPTS"))
	var netClient = &http.Client{
		Timeout: time.Second * time.Duration(timeOutSec),
	}

	var respBody []byte
	var statusCode int

	// Store body to repopulate body for retry requests
	bodyData := []byte{}
	if req.Method == "POST" && req.Body != nil {
		bodyData, _ = ioutil.ReadAll(req.Body)
		req.Body.Close()
	}

	err := RetryCall(attempts, time.Second, func() error {
		// Need to populate body every time as the body becomes empty after first request.
		if req.Method == "POST" {
			req.Body = ioutil.NopCloser(bytes.NewReader(bodyData))
		}

		resp, err := netClient.Do(req)
		if err != nil {
			return fmt.Errorf("server error: %v", err)
		}
		defer resp.Body.Close()

		statusCode = resp.StatusCode
		respBody, err = ioutil.ReadAll(resp.Body)

		switch {
		case statusCode >= 500:
			return fmt.Errorf("server error: %v", statusCode)
		case statusCode >= 400:
			return stop{fmt.Errorf("client error: %v", statusCode)}
		default:
			return nil
		}
	})

	// timeout and no such host error are StatusInternalServerError
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	return respBody, statusCode, err
}

func RetryCall(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return RetryCall(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}
