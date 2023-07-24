package common

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "bucketctl"), nil
}

func CreateDirIfNotExists(dir string, perm os.FileMode) error {
	baseDir := path.Dir(dir)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, perm)
}

func FileNotExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true
		}
	}
	return false
}

func CreateFile(file string, perm os.FileMode) error {
	var fileHandle, err = os.Create(file)
	defer fileHandle.Close()
	if err != nil {
		return err
	}
	if err := os.Chmod(file, perm); err != nil {
		return err
	}
	return nil
}

func CreateFileIfNotExists(file string, perm os.FileMode) error {
	if FileNotExists(file) {
		if err := CreateFile(file, perm); err != nil {
			return err
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

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func RemoveFile(filename string) error {
	return os.Remove(filename)
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
