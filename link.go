package main

import (
	"errors"
	"log"
	"os"
	"path"
)

type links map[string]string

func (l *links) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux map[string]string

	if err := unmarshal(&aux); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	home := os.Getenv("HOME")

	for target, dest := range aux {
		target = path.Join(cwd, os.ExpandEnv(target))
		dest = path.Join(home, os.ExpandEnv(dest))

		aux[target] = dest
	}

	*l = links(aux)

	return nil
}

func isSymlink(link string) (bool, error) {
	stat, err := os.Lstat(link)
	if os.IsNotExist(err) {
		return false, err
	}

	if stat.Mode()&os.ModeSymlink == 0 {
		return false, nil
	}

	return true, nil
}

func removeLink(link string) error {
	ok, err := isSymlink(link)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("File is not a symlink")
	}

	if err := os.Remove(link); err != nil {
		return err
	}

	log.Printf("Symlink removed: %s", link)

	return nil
}

func removeLinks() {

}

func clobberPath(path string) error {
	stat, err := os.Lstat(path)
	if os.IsNotExist(err) {
		// we return here because there is no file to clean
		return nil
	}

	// remove the file/dir if it is not a symlink
	if stat.Mode()&os.ModeSymlink == 0 {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	} else {
		err := removeLink(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func createParents(link string) error {
	parentDir := path.Dir(link)

	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		log.Printf("Creating path: %s", parentDir)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (t *task) createLinks(flags *provisionFlags) error {
	for target, dest := range t.Links {

		if flags.parents {
			err := createParents(dest)
			if err != nil {
				return err
			}
		}

		if flags.clobber {
			err := clobberPath(dest)
			if err != nil {
				return err
			}
		}

		err := os.Symlink(target, dest)
		if err != nil {
			return err
		}

		log.Printf("Symlink created: %s", dest)
	}

	return nil
}
