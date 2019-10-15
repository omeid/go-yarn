package catalog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	yarn "github.com/omeid/go-yarn"
)

var (
	// ErrUnknownCommentTag is return when a line that is expected to be
	// a starting line starts with an unsupported comment type.
	// i.e: it is not `--`, `#`, `//`, or `/*`.
	ErrUnknownCommentTag = errors.New("Uknown comment pair")

	// ErrUnexpectedEnd is returned when a catalog is malformed.
	ErrUnexpectedEnd = errors.New("Unexpected end")
)

// MustParse is like Parse, but panics if the named string
// is not found in the provied Yarn or fails to parse
// the catalog.
func MustParse(y yarn.Yarn, name string) yarn.Yarn {
	yarn, err := Parse(y, name)
	if err != nil {
		panic(err)
	}

	return yarn
}

const nl = "\n"

const (
	commentAdaLike  = "--"
	commentBashLike = "##"
	commentCLike    = "//"

	commentCBlockStart = "/*"
	commentCBlockEnd   = "*/"
)

func endTag(start string) (string, error) {
	switch start {
	case commentAdaLike,
		commentBashLike,
		commentCLike:
		return "", nil // requires no endtag.
	case commentCBlockStart:
		return commentCBlockEnd, nil
	default:
		return "", ErrUnknownCommentTag
	}
}

// Parse the value of catalog from the provided yarn.
func Parse(y yarn.Yarn, catalog string) (yarn.Yarn, error) {
	source, ok := y.Get(catalog)
	if !ok {
		return nil, fmt.Errorf(yarn.MissingYarn, catalog)
	}

	return parse(source)
}

// ParseString is like Parse, but accept a string as source.
func ParseString(source string) (yarn.Yarn, error) {
	return parse(source)
}

// ParseFile reads the content of the file and parses into a Yarn.
func ParseFile(filename string) (yarn.Yarn, error) {
	bs, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return parse(string(bs))
}

// Parse loads the value of name from y into a new Yarn.
func parse(source string) (yarn.Yarn, error) {

	files := map[string]string{}
	var (
		// per catalog.
		start *regexp.Regexp
		end   *regexp.Regexp

		// per entry.
		entryName    string
		entryContent string
		err          error
	)

	entryAdd := func() {
		files[entryName] = entryContent
		entryName = ""
		entryContent = ""
	}

	entryReadName := func(i int, line string) error {

		// only once, we only support one comment type per catalog.
		if start == nil {
			var err error
			start, end, err = getDelims(line)
			if err != nil {
				return err
			}
		}

		segs := start.FindStringSubmatch(line)

		if len(segs) < 2 {
			return fmt.Errorf("Invalid start tag at line %v", i)
		}

		entryName = segs[1]
		return nil
	}

	entryEndline := func(i int, line string) (bool, error) {

		segs := end.FindStringSubmatch(line)
		if len(segs) == 0 {
			return false, nil
		}

		if len(segs) < 2 {
			return true, fmt.Errorf("Invalid end tag at line %v", i)
		}

		if segs[1] != entryName {
			return true, fmt.Errorf("Invalid end tag name %v expected %v at line %v", segs[1], entryName, i)
		}

		return true, nil
	}

	for i, line := range strings.Split(source, nl) {

		if entryName != "" {
			end, err := entryEndline(i, line)

			if err != nil {
				return nil, err
			}

			if end {
				entryAdd()
				continue
			}

			// Don't add newline to the start of the content.
			if entryContent != "" {
				line = nl + line
			}

			entryContent += line
			continue
		}

		// whitespace only.
		if len(strings.Fields(line)) == 0 {
			continue
		}

		// then it must be starting tag.
		err = entryReadName(i, line)
		if err != nil {
			return nil, err
		}
	}

	if entryName != "" {
		// Unexpected end.
		return nil, ErrUnexpectedEnd
	}

	return yarn.NewFromMap(files), nil
}

func getDelims(line string) (*regexp.Regexp, *regexp.Regexp, error) {
	if len(line) < 2 {
		return nil, nil, ErrUnknownCommentTag
	}

	labels := []string{"start", "end"}

	start := line[:2]
	end, err := endTag(start)

	if err != nil {
		return nil, nil, err
	}

	regs, err := makeRegexs(labels, start, end)
	if err != nil {
		return nil, nil, err
	}

	return regs[0], regs[1], nil
}

func makeRegexs(labels []string, start, end string) ([]*regexp.Regexp, error) {
	regs := make([]*regexp.Regexp, len(labels))

	var err error

	start = regexp.QuoteMeta(start)
	end = regexp.QuoteMeta(end)
	for i, label := range labels {

		reg := start + `\s+` + label + `\s?:\s*(.*)?`
		if end != "" {
			reg = reg + `\s+` + end + `\s*$`
		} else {
			reg = reg + `\s*$`
		}

		regs[i], err = regexp.Compile(reg)
		if err != nil {
			return nil, err
		}
	}
	return regs, nil
}
