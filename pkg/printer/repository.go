package printer

import (
	"bucketctl/pkg/api/bitbucket/types"
	"bucketctl/pkg/common"
	"strconv"
)

func PrettyFormatRepositories(reposMap map[string]*types.Repository) [][]string {
	var data [][]string
	data = append(data, []string{"ID", "Slug", "State", "Public", "Archived"})

	repos := common.GetLexicallySortedKeys(reposMap)
	for _, slug := range repos {
		row := []string{strconv.Itoa(reposMap[slug].Id), slug, reposMap[slug].StatusMessage, strconv.FormatBool(reposMap[slug].Public), strconv.FormatBool(reposMap[slug].Archived)}
		data = append(data, row)
	}

	return data
}
