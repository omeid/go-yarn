// Filesystem mapped key-string store. Ideal for embedding code like sql.
package yarn

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

//Like New, but panics on error.
func Must(fs http.FileSystem, pattern string) *Yarn {
	y, e := New(fs, pattern)
	if e != nil {
		panic(e)
	}
	return y
}

//Creates a new Yarn from provided filesystem's files that match the pattern,
func New(fs http.FileSystem, pattern string) (*Yarn, error) {

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

	yarn := &Yarn{make(map[string]string)}
	for _, file := range files {
		name := file.Name()
		//the pattern is already checked in the start and filepath.Match doesn't return fileysstem errors, so ignoring the error.
		if ok, _ := filepath.Match(pattern, name); !ok {
			continue
		}
		file, err := fs.Open(name)
		if err != nil {
			return yarn, err
		}

		content, err := ioutil.ReadAll(file)
		//conntent, err := ioutil.Readall(file)
		if err != nil {
			return yarn, err
		}
		yarn.strings[name] = string(content)

	}
	return yarn, nil
}

type Yarn struct {
	strings map[string]string
}

func (y *Yarn) Has(strings ...string) error {
	var (
		s  string
		ok bool
	)
	for _, s = range strings {
		if _, ok = y.strings[s]; !ok {
			return fmt.Errorf(" Missing %s.", s)
		}
	}
	return nil
}

//Returns a loaded file's contents as string and if it exists by filename.
func (y *Yarn) Get(key string) (string, bool) {
	content, ok := y.strings[key]
	return content, ok
}

//Returns a loaded file's contents as string, it panics if file doesn't exist.
const missingYarn = "Yarn missing %s"
func (y *Yarn) Must(key string) string {
	content, ok := y.strings[key]
	if !ok {
		panic(fmt.Sprintf(missingYarn, key))
	}
	return content
}
