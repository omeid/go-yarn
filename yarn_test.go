package yarn

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
)

var (
	testyarns map[string]string
	goodkeys  []string
	badkeys   []string
	testdata  Yarn
)

func TestMain(m *testing.M) {
	testyarns = map[string]string{
		"insert.sql":    "INSERT INTO users (id, name, email)\nVALUES ($1, $2, $3)\n",
		"query_all.sql": "SELECT\nid,\nname,\nFROM users\n",
		"web/pages.js":  "'strict'\n\n;(function () {\n}())\n",
		"web/test.css":  "#app {\n}\n",
	}

	goodkeys = []string{"insert.sql", "query_all.sql"}

	badkeys = []string{"something.json", "random", "sql", "insert", "nope", "none"}
	testdata = Must(http.Dir("testdata"), "*.sql", "web/*")

	os.Exit(m.Run())
}

func TestMustHave(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("Didn't panic for bad keys: %s", badkeys)
		}
		if r == fmt.Sprintf(MissingYarn, badkeys) {
			return
		}
		panic(r)
	}()
	testdata.MustHave(badkeys...)
	//TODO: Handle this properly.
	testdata.MustHave(goodkeys...)
}

func TestHas(t *testing.T) {

	err := testdata.Has(goodkeys...)
	if err != nil {
		t.Fatal(err)
	}

	if err := testdata.Has(badkeys...); err == nil {
		t.Fatal("Expected error. Got nothing.")
	}
}

func TestMust(t *testing.T) {
	for name, testcontent := range testyarns {
		func(name, testcontent string) {
			defer func() {
				r := recover()
				if r != nil {
					if r == fmt.Sprintf(MissingYarn, name) {
						t.Fatalf("Missing using MUST %s", name)
						return
					}
					panic(r)
				}
			}()
			content := testdata.Must(name)
			if content != testcontent {
				t.Fatalf("For %s:\nExpected:\n`%s`\nGot:\n`%s`\n", name, testcontent, content)
			}
		}(name, testcontent)
	}

	for _, name := range badkeys {
		func(name string) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("Must didn't panic for unexpected `%s` key.", name)
				}
				if r != fmt.Sprintf(MissingYarn, name) {
					panic(r)
				}
			}()
			testdata.Must(name)

		}(name)
	}
}

func TestGet(t *testing.T) {
	for name, testcontent := range testyarns {
		content, ok := testdata.Get(name)
		if !ok {
			t.Fatalf("Missing %s", name)
		}
		if content != testcontent {
			t.Fatalf("For %s:\nExpected:\n`%s`\nGot:\n`%s`\n", name, testcontent, content)
		}
	}

	for _, name := range badkeys {
		if _, ok := testdata.Get(name); ok {
			t.Fatalf("Got OK for unexpected `%s` key.", name)
		}
	}
}

func TestSub(t *testing.T) {

	web := testdata.Sub("web")
	for name, testcontent := range testyarns {
		if !strings.HasPrefix(name, "web/") {
			continue
		}

		name = strings.TrimPrefix(name, "web/")
		content, ok := web.Get(name)
		if !ok {
			t.Fatalf("Missing %s", name)
		}
		if content != testcontent {
			t.Fatalf("For %s:\nExpected:\n`%s`\nGot:\n`%s`\n", name, testcontent, content)
		}

	}
}
