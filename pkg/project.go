package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Project struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type Projects struct {
	Response
	Values []Project `json:"values"`
}

func GetProjects(baseUrl string, limit int) Projects {
	url := fmt.Sprintf("%s/rest/api/1.0/projects/?limit=%d", baseUrl, limit)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var result Projects
	if err := json.Unmarshal(body, &result); err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}

	return result
}

func PrintProjects(projects []Project) {
	var data [][]string

	data = append(data, []string{"ID", "Key", "Name", "Description"})

	for _, proj := range projects {
		row := []string{strconv.Itoa(proj.Id), proj.Key, proj.Name, proj.Description}
		data = append(data, row)
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
