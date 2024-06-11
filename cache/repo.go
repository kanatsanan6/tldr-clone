package cache

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/kanatsanan6/tldr/utils"
)

type Repository struct {
	remoteUrl string
	dir       string
}

func NewRepository(remoteUrl string, ttl time.Duration) (*Repository, error) {
	dir, err := cacheDir()
	if err != nil {
		return nil, fmt.Errorf("error while getting cache directory: %s", err)
	}
	repo := &Repository{remoteUrl: remoteUrl, dir: dir}
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Println("cache does not exist")

		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return nil, fmt.Errorf("error while creating cache directory: %s", err)
			}
		} else {
			return nil, fmt.Errorf("error while creating cache directory: %s", err)
		}

		if err := repo.create(); err != nil {
			return nil, fmt.Errorf("error while downloading the information: %s", err)
		}
	} else {
		fmt.Println("caching directory already exists")

		if info.ModTime().Before(time.Now().Add(-ttl)) {
			repo.reload()
		}
	}
	return repo, nil
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

func (r *Repository) create() error {
	err := utils.Download(r.remoteUrl, filepath.Join(r.dir, "tldr.zip"))
	if err != nil {
		return err
	}
	err = utils.Unzip(filepath.Join(r.dir, "tldr.zip"), r.dir)
	return err
}

func (r *Repository) reload() error {
	err := os.RemoveAll(r.dir)
	if err != nil {
		return err
	}
	return r.create()
}
