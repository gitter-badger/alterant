package main

import (
	"errors"
	"log"
	"os"
	"path"
)

type linkTarget struct {
	Value string
}

type linkDestination struct {
	Value string
}

// TODO: clean this up
func (l *linkDestination) isSymlink() (bool, error) {
	stat, err := os.Lstat(l.Value)
	if os.IsNotExist(err) {
		return false, err
	}

	if stat.Mode()&os.ModeSymlink == 0 {
		return false, nil
	}

	return true, nil
}

// TODO: clean this up
func (l *linkDestination) removeLink() error {
	ok, err := l.isSymlink()

	if err != nil {
		return err
	}

	if !ok {
		return errors.New("File is not a symlink")
	}

	if err := os.Remove(l.Value); err != nil {
		return err
	}

	log.Printf("Symlink removed: %s", l.Value)

	return nil
}

func (l *linkDestination) clobberPath() error {
	stat, err := os.Lstat(l.Value)
	if os.IsNotExist(err) {
		// we return here because there is no file to clean
		return nil
	}

	// remove the file/dir if it is not a symlink
	if stat.Mode()&os.ModeSymlink == 0 {
		err := os.RemoveAll(l.Value)
		if err != nil {
			return err
		}
	} else {
		err := l.removeLink()
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
			err := createParents(dest.Value)
			if err != nil {
				return err
			}
		}

		if flags.clobber {
			err := dest.clobberPath()
			if err != nil {
				return err
			}
		}

		err := os.Symlink(target.Value, dest.Value)
		if err != nil {
			return err
		}

		log.Printf("Symlink created: %s", dest.Value)
	}

	return nil
}
