package cache

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/kanatsanan6/tldr/utils"
)

type Repository struct {
	dir string
}

func NewRepository(remoteUrl string) (*Repository, error) {
	dir, err := cacheDir()
	if err != nil {
		return nil, fmt.Errorf("error while getting cache directory: %s", err)
	}

	err = utils.CreateIfNotExist(dir)
	if err != nil {
		return nil, fmt.Errorf("error while creating cache directory: %s", err)
	}

	err = utils.Download(remoteUrl, filepath.Join(dir, "tldr.zip"))
	if err != nil {
		return nil, fmt.Errorf("error while downloading cache file: %s", err)
	}

	err = utils.Unzip(filepath.Join(dir, "tldr.zip"), dir)
	if err != nil {
		return nil, fmt.Errorf("error while extracting cache file: %s", err)
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
