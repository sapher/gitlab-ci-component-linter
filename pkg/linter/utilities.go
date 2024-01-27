package linter

import "os"

func IsFileExist(path string) bool {
	stat, err := os.Stat(path)
	return os.IsExist(err) && !stat.IsDir()
}

func IsDirExist(path string) bool {
	stat, err := os.Stat(path)
	return os.IsExist(err) && stat.IsDir()
}
