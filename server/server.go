package server

import (
	"os"
	"path/filepath"

	"github.com/davidlouie/mpgo/server/subsonic"
)

func GetFiles(root string) ([]string, error) {
	ext := "*.mp3"
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(ext, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func Init() {
	subsonic.Init()
}
