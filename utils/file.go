package utils

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func CreateDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return err
}

func Download(remoteUrl, des string) error {
	out, err := os.Create(des)
	if err != nil {
		return err
	}

	res, err := http.Get(remoteUrl)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func Unzip(zipPath, des string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	extractAndWrite := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// TODO: check for ZipSlip before creating the file
		path := filepath.Join(des, f.Name)
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
			return err
		}
	}
	return nil
}
