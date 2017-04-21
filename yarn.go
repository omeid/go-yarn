// Package yarn is a filesystem mapped key-string store. Ideal for embedding code like sql.
package yarn

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)

// Must is like New, but panics on error.
func Must(fs http.FileSystem, pattern string) Yarn {
	y, e := New(fs, pattern)
	if e != nil {
		panic(e)
	}
	return y
}

// Yarn is script store.
type Yarn interface {
	// Has checks if a file for the provided list of keys exists, if not, returns an error.
	Has(strings ...string) error
	// MustHave is like Has but panics on missing keys.
	MustHave(strings ...string)
	// Get returns a loaded file's contents as string and if it exists by filename.
	Get(key string) (string, bool)
	// Must returns a loaded file's contents as string, it panics if file doesn't exist.
	Must(key string) string

	Sub(dir string) Yarn
}

// New creates a new Yarn from provided filesystem's files that match the pattern,
func New(fs http.FileSystem, pattern string) (Yarn, error) {

	//Check the pattern.
	_, err := filepath.Match(pattern, "")
	if err != nil {
		return nil, err
	}
	dir, err := fs.Open("/")
	if err != nil {
		return nil, err
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	yarn := newYarn()
	for _, file := range files {
		name := file.Name()

		ok, err := filepath.Match(pattern, name)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		file, err := fs.Open(name)
		if err != nil {
			return yarn, err
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			return yarn, err
		}
		yarn.add(name, string(content))

	}
	return yarn, nil
}
