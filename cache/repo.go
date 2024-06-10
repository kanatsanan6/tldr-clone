package cache

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

type Repository struct {
	dir string
}

func NewRepository(remoteUrl string) (*Repository, error) {
	dir, err := cacheDir()
	if err != nil {
		return nil, fmt.Errorf("error while getting cache directory: %s", err)
	}

	// check if caching dir exists
	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return nil, fmt.Errorf("error while creating cache directory: %s", err)
			}
		} else {
			return nil, fmt.Errorf("somethings wrong with caching dir: %s", err)
		}
	}

	out, err := os.Create(filepath.Join(dir, "tldr.zip"))
	if err != nil {
		return nil, fmt.Errorf("error while creating cache file: %s", err)
	}

	resp, err := http.Get(remoteUrl)
	if err != nil {
		return nil, fmt.Errorf("error while downloading the file: %s", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while copying the file: %s", err)
	}

	// unzip
	reader, err := zip.OpenReader(out.Name())
	if err != nil {
		return nil, fmt.Errorf("error while opening the file: %s", err)
	}
	defer reader.Close()

	extractAndWrite := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// TODO: check for ZipSlip before creating the file
		path := filepath.Join(dir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for _, f := range reader.File {
		err := extractAndWrite(f)
		if err != nil {
			return nil, fmt.Errorf("error while extracting the file: %s", err)
		}
	}

	return &Repository{dir: dir}, nil
}

func cacheDir() (string, error) {
	// try to use XDG_CACHE_HOME for caching
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir != "" {
		return filepath.Join(cacheDir, "tldr_clone"), nil
	}

	// if there is no XDG_CACHE_HOME, use user's home dir instead
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error while getting current user info: %s", err)
	}
	return filepath.Join(user.HomeDir, ".cache", "tldr_clone"), nil
}
