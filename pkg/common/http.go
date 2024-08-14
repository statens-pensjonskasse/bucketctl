package common

import (
	"encoding/json"
	"fmt"
	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HttpRequest(method string, url string, payload io.Reader, token string, params ...map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, payload)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if params != nil && len(params) > 0 {
		q := req.URL.Query()
		for key, val := range params[0] {
			if strings.HasPrefix(key, "Header ") {
				headerKey := strings.TrimPrefix(key, "Header ")
				req.Header.Set(headerKey, val)
			} else {
				q.Add(key, val)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	// Http-status 429 tyder på at vi sender for mange kall. Venter og prøver igjen.
	if resp.StatusCode == 429 {
		wait, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		time.Sleep(time.Duration(wait) * time.Second)
		return HttpRequest(method, url, payload, token, params...)
	}

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		var errorResp types.Error
		if err := json.Unmarshal(bodyBytes, &errorResp); err != nil {
			return resp, fmt.Errorf("http status %d for %s-call to %s: %s", resp.StatusCode, method, url, string(bodyBytes))
		}
		return resp, fmt.Errorf("http status %d for %s-call to %s: %s", resp.StatusCode, method, url, errorResp.Errors)
	}

	return resp, nil
}

func GetRequest(url string, token string) (*http.Response, error) {
	return HttpRequest("GET", url, nil, token)
}

func DeleteRequest(url string, token string, params map[string]string) (*http.Response, error) {
	return HttpRequest("DELETE", url, nil, token, params)
}

func PostRequest(url string, token string, payload io.Reader, params map[string]string) (*http.Response, error) {
	return HttpRequest("POST", url, payload, token, params)
}

func PutRequest(url string, token string, payload io.Reader, params map[string]string) (*http.Response, error) {
	return HttpRequest("PUT", url, payload, token, params)
}

func GetRequestBody(url string, token string) ([]byte, error) {
	resp, err := GetRequest(url, token)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
