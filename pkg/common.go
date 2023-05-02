package pkg

import (
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"io/fs"
	"net/http"
	"os"
)

type BitbucketResponse struct {
	Size          int    `json:"size"`
	Limit         int    `json:"limit"`
	IsLastPage    bool   `json:"isLastPage"`
	Start         int    `json:"start"`
	NextPageStart int    `json:"nextPageStart"`
	Values        []byte `json:"values"`
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

func HttpRequest(method string, url string, body io.Reader, token string) (*http.Response, error) {
	client := http.Client{}
	req, _ := http.NewRequest(method, url, body)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 300 {
		return resp, fmt.Errorf("http status %d for %s-call to %s", resp.StatusCode, method, url)
	}

	return resp, nil
}

func GetRequest(url string, token string) (*http.Response, error) {
	return HttpRequest("GET", url, nil, token)
}

func GetRequestBody(url string, token string) ([]byte, error) {
	resp, err := GetRequest(url, token)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
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
