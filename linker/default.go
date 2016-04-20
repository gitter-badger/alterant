package linker

// TODO: add a command to remove cruft and broken links from the home folder
// that have been managed by alterant

import (
	"errors"
	"os"
	"path"

	"github.com/autonomy/alterant/link"
	"github.com/autonomy/alterant/logger"
)

//DefaultLinker is a basic symlink manager and is the default
type DefaultLinker struct {
	logger  *logWrapper.LogWrapper
	enabled bool
	parents bool
	clobber bool
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

func (dl *DefaultLinker) removeLink(link string) error {
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

	dl.logger.Info(2, "Symlink removed: %s", link)

	return nil
}

func (dl *DefaultLinker) clobberPath(path string) error {
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
		err := dl.removeLink(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dl *DefaultLinker) createParents(link string) error {
	parentDir := path.Dir(link)

	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		dl.logger.Info(2, "Creating path: %s", parentDir)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// RemoveLinks removes symlinks
func (dl *DefaultLinker) RemoveLinks(links []*link.Link) error {
	for _, link := range links {
		err := dl.removeLink(string(link.Destination))
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateLinks creates symlinks
func (dl *DefaultLinker) CreateLinks(links []*link.Link) error {

	// do not create links if not enabled
	if !dl.enabled {
		return nil
	}

	for _, link := range links {
		if dl.parents {
			err := dl.createParents(string(link.Destination))
			if err != nil {
				return err
			}
		}

		if dl.clobber {
			err := dl.clobberPath(string(link.Destination))
			if err != nil {
				return err
			}
		}

		// TODO: validate symlinks
		err := os.Symlink(string(link.Target), string(link.Destination))
		if err != nil {
			return err
		}

		dl.logger.Info(2, "Symlink created: %s -> %s", link.Destination, link.Target)
	}

	return nil
}

// NewDefaultLinker returns an instance of `DefaultLinker`
func NewDefaultLinker(enabled bool, parents bool, clobber bool, logger *logWrapper.LogWrapper) *DefaultLinker {
	return &DefaultLinker{
		logger:  logger,
		enabled: enabled,
		parents: parents,
		clobber: clobber,
	}
}
