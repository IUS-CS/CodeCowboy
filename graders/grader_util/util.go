package util

import (
	cp "github.com/otiai10/copy"
	"os"
	"path"
)

func CopyExtras(from string, to string) error {
	d, err := os.ReadDir(from)
	if err != nil {
		return err
	}
	for _, entry := range d {
		err = cp.Copy(path.Join(from, entry.Name()), to)
		if err != nil {
			return err
		}
	}
	return nil
}
