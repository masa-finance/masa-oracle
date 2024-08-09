package workers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	GET   = "GET"
	POST  = "POST"
	PUT   = "PUT"
	PATCH = "PATCH"
)

var transport *http.Transport

func getHTTPClient() http.Client {
	if transport == nil {
		transport = &http.Transport{
			MaxIdleConnsPerHost: 0,
			TLSHandshakeTimeout: 5 * time.Second,
		}
	}
	var netClient = http.Client{
		Timeout:   time.Second * 300,
		Transport: transport,
	}
	return netClient
}

func Get(url string, headers map[string]string) ([]byte, error) {
	return call(url, GET, bytes.NewBuffer(make([]byte, 0)), headers)
}

func Post(url string, rawJSON json.RawMessage, headers map[string]string) ([]byte, error) {
	return call(url, POST, bytes.NewBuffer(rawJSON), headers)
}

func Put(url string, rawJSON json.RawMessage, headers map[string]string) ([]byte, error) {
	return call(url, PUT, bytes.NewBuffer(rawJSON), headers)
}

func Patch(url string, rawJSON json.RawMessage, headers map[string]string) ([]byte, error) {
	return call(url, PATCH, bytes.NewBuffer(rawJSON), headers)
}

func call(url, method string, buffer *bytes.Buffer, headers map[string]string) ([]byte, error) {
	client := getHTTPClient()
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, err
	}
	contentSet := false
	first := true
	if headers != nil {
		for key, value := range headers {
			if key == "Content-Type" {
				contentSet = true
			}
			if first {
				first = false
				req.Header.Set(key, value)
			} else {
				req.Header.Add(key, value)
			}
		}
		// default to application/json if content type is not specified
		if !contentSet {
			if first {
				req.Header.Set("Content-Type", "application/json")
			} else {
				req.Header.Add("Content-Type", "application/json")
			}
			req.Header.Add("Accept-Encoding", "gzip, deflate, br")
		}
	}
	resp, err := client.Do(req)
	if resp == nil {
		if err != nil {
			return nil, err
		} else {
			return nil, errors.New("http response was null")
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Error closing response body: ", err)
		}
	}(resp.Body)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	if resp.StatusCode > 299 {
		return body, errors.New(resp.Status)
	}
	return body, err
}
