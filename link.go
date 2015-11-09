package main

import (
	"log"
	"os"
	"path"
)

func removeLink(link string) error {
	stat, err := os.Lstat(link)
	if os.IsNotExist(err) || stat.Mode()&os.ModeSymlink == 0 {
		return nil
	}

	if err := os.Remove(link); err != nil {
		return err
	}

	log.Printf("Symlink removed: %s", link)

	return nil
}

func cleanPath(link string) error {
	stat, err := os.Lstat(link)
	if os.IsNotExist(err) {
		// we return here because there is no file to clean
		return nil
	}

	// remove the file/dir if it is not a symlink
	if stat.Mode()&os.ModeSymlink == 0 {
		err := os.RemoveAll(link)
		if err != nil {
			return err
		}
	} else {
		err := removeLink(link)
		if err != nil {
			return err
		}
	}

	return nil
}

func createPath(destination string) error {
	parentDir := path.Dir(destination)

	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		log.Printf("Creating path: %s", parentDir)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (t *task) createLinks(cwd string) error {
	for original, link := range t.Links {
		source := path.Join(cwd, original)
		destination := path.Join(os.ExpandEnv(t.Destination), link)

		err := createPath(destination)
		if err != nil {
			return err
		}

		err = cleanPath(destination)
		if err != nil {
			return err
		}

		err = os.Symlink(os.ExpandEnv(source), os.ExpandEnv(destination))
		if err != nil {
			return err
		}

		log.Printf("Symlink created: %s", destination)
	}

	return nil
}
