package pkg

import (
	"errors"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
)

type Response struct {
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
