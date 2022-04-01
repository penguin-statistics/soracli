package filepath

import (
	"os"
	"path"

	"github.com/penguin-statistics/soracli/internal/consts"
)

func UnderDataDir(p string) string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	file := path.Join(userHomeDir, consts.DataDir, p)
	dir := path.Dir(file)
	// if dir does not exist, create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(err)
		}
	}

	return file
}
