package printer

import (
	"strconv"

	"git.spk.no/infra/bucketctl/pkg/api/bitbucket/types"
	"git.spk.no/infra/bucketctl/pkg/common"
)

func PrettyFormatProjects(projectsMap map[string]*types.Project) [][]string {
	var data [][]string
	data = append(data, []string{"ID", "Project Key", "Type", "Description"})

	projects := common.GetLexicallySortedKeys(projectsMap)
	for _, key := range projects {
		row := []string{strconv.Itoa(projectsMap[key].Id), key, projectsMap[key].Name, projectsMap[key].Description}
		data = append(data, row)
	}

	return data
}
