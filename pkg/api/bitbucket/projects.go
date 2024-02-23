package bitbucket

import (
	"bucketctl/pkg/api/bitbucket/types"
	"bucketctl/pkg/common"
	"bucketctl/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func GetProjects(baseUrl string, limit int, token string) (map[string]*types.Project, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects?limit=%d", baseUrl, limit)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var projectsResponse types.ProjectsResponse
	if err := json.Unmarshal(body, &projectsResponse); err != nil {
		return nil, err
	}

	if !projectsResponse.IsLastPage {
		logger.Warn("Not all projects fetched, try with a higher limit")
	}

	projects := make(map[string]*types.Project)
	for _, p := range projectsResponse.Values {
		projects[p.Key] = &types.Project{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Public:      p.Public,
			Type:        p.Type,
			Links:       p.Links,
		}
	}

	return projects, nil
}

func GetProject(baseUrl string, projectKey string, token string) (*types.Project, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s", baseUrl, projectKey)

	body, err := common.GetRequestBody(url, token)
	if err != nil {
		return nil, err
	}

	var project types.Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func UpdateProject(baseUrl string, token string, project *types.Project) (*types.Project, error) {
	url := fmt.Sprintf("%s/rest/api/latest/projects/%s", baseUrl, project.Key)

	payload, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	resp, err := common.PutRequest(url, token, bytes.NewReader(payload), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var p types.Project
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
