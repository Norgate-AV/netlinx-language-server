package queries

import (
	"embed"
)

//go:embed *.scm
var QueryFS embed.FS

func GetQuery(path string) (string, error) {
	data, err := QueryFS.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
