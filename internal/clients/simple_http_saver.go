package clients

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

type SimpleHttpSaver struct {
}

func (s SimpleHttpSaver) Save(url, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create dir: %s error: %w", dir, err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot get data from url: %s error: %w", url, err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("cannot create file: %s error: %w", filePath, err)
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	if _, err = io.Copy(io.MultiWriter(file, bar), resp.Body); err != nil {
		return fmt.Errorf("cannot save file: %s error: %w", filePath, err)
	}
	return nil
}
