package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

type BitbucketResponse struct {
	Size          int    `json:"size"`
	Limit         int    `json:"limit"`
	IsLastPage    bool   `json:"isLastPage"`
	Start         int    `json:"start"`
	NextPageStart int    `json:"nextPageStart"`
	Values        []byte `json:"values"`
}

type Group struct {
	Name string `json:"name" yaml:"name"`
}

type User struct {
	Name         string `json:"name" yaml:"name"`
	EmailAddress string `json:"emailAddress" yaml:"emailAddress"`
	Active       bool   `json:"active" yaml:"active"`
	DisplayName  string `json:"displayName" yaml:"displayName"`
	Id           int    `json:"id" yaml:"id"`
	Slug         string `json:"slug" yaml:"slug"`
	Type         string `json:"type" yaml:"type"`
}

func CreateFileIfNotExists(file string) {
	if _, err := os.Stat(file); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			var file, err = os.Create(file)
			defer file.Close()
			if err != nil {
				pterm.Error.Println("Error creating config file:", file)
				os.Exit(1)
			}
		}
	}
}

func HttpRequest(method string, url string, body io.Reader, token string, params ...map[string]string) (*http.Response, error) {
	client := http.Client{}
	req, _ := http.NewRequest(method, url, body)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	if params != nil && len(params) > 0 {
		q := req.URL.Query()
		for key, val := range params[0] {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf("http status %d for %s-call to %s: %s", resp.StatusCode, method, url, string(bodyBytes))
	}

	return resp, nil
}

func GetRequest(url string, token string) (*http.Response, error) {
	return HttpRequest("GET", url, nil, token)
}

func DeleteRequest(url string, token string, params map[string]string) (*http.Response, error) {
	return HttpRequest("DELETE", url, nil, token, params)
}

func PutRequest(url string, token string, params map[string]string) (*http.Response, error) {
	return HttpRequest("PUT", url, nil, token, params)
}

func GetRequestBody(url string, token string) ([]byte, error) {
	resp, err := GetRequest(url, token)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func ReadConfigFile[T interface{}](filename string, obj T) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(filename, ".yaml") {
		if err := yaml.Unmarshal(file, &obj); err != nil {
			return err
		}
	} else if strings.HasSuffix(filename, ".json") {
		if err := json.Unmarshal(file, &obj); err != nil {
			return err
		}
	} else {
		return errors.New("forventet fil med enten .yaml eller .json ending")
	}

	return nil
}
