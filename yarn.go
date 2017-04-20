// Package yarn is a filesystem mapped key-string store. Ideal for embedding code like sql.
package yarn

import (
	"fmt"
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

	yarn := &yarn{make(map[string]string)}
	for _, file := range files {
		name := file.Name()
		//the pattern is already checked in the start so we ignore the error.
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
		yarn.strings[name] = string(content)

	}
	return yarn, nil
}

type yarn struct {
	strings map[string]string
}

const missingYarn = "Missing %s"

// Has checks if a file for the provided list of keys exists, if not, returns an error.
func (y *yarn) Has(strings ...string) error {
	var (
		s       string
		ok      bool
		missing []string
	)

	for _, s = range strings {
		if _, ok = y.strings[s]; !ok {
			missing = append(missing, s)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf(missingYarn, missing)
	}
	return nil
}

// MustHave is like Has but panics on missing keys.
func (y *yarn) MustHave(strings ...string) {
	err := y.Has(strings...)
	if err != nil {
		panic(err.Error())
	}
}

// Get returns a loaded file's contents as string and if it exists by filename.
func (y *yarn) Get(key string) (string, bool) {
	content, ok := y.strings[key]
	return content, ok
}

// Must returns a loaded file's contents as string, it panics if file doesn't exist.
func (y *yarn) Must(key string) string {
	content, ok := y.strings[key]
	if !ok {
		panic(fmt.Sprintf(missingYarn, key))
	}
	return content
}
