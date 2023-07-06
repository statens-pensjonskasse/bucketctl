package pkg

import (
	"bucketctl/pkg/types"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func CreateDirIfNotExists(dir string, perm os.FileMode) error {
	baseDir := path.Dir(dir)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, perm)
}

func CreateFileIfNotExists(file string, perm os.FileMode) error {
	if _, err := os.Stat(file); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			var fileHandle, err = os.Create(file)
			defer fileHandle.Close()
			if err != nil {
				return err
			}
			if err := os.Chmod(file, perm); err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckFilePermission(file string, perm os.FileMode) error {
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}
	if stat.Mode() != perm {
		return errors.New("Unexpected file permission '" + stat.Mode().String() + "' for file '" + file + "', expected '" + perm.String() + "'")
	}
	return nil
}

func HttpRequest(method string, url string, payload io.Reader, token string, params ...map[string]string) (*http.Response, error) {
	// TODO: DON'T SKIP TLS VERIFICATION!!!
	// Temporary workaround for Cert-issues on Mac
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
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

func GetLexicallySortedKeys[T any](stringMap map[string]T) []string {
	keys := make([]string, 0, len(stringMap))
	for k := range stringMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func SlicesContainsSameElements[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	diff := make(map[T]int, len(a))
	for _, i := range a {
		// Tell antall ganger verdi dukker opp
		diff[i]++
	}
	for _, j := range b {
		if _, exists := diff[j]; !exists {
			return false
		}
		diff[j]--
		if diff[j] == 0 {
			// Slett dersom vi har funnet elementet nok ganger
			delete(diff, j)
		}
	}
	return len(diff) == 0
}
