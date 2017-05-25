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

	goodkeys = []string{"insert.sql", "query_all.sql", "web/pages.js", "web/test.css"}

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

		if strings.HasPrefix(name, "web/") {

			name = strings.TrimPrefix(name, "web/")
			content, ok := web.Get(name)
			if !ok {
				t.Fatalf("Missing %s", name)
			}
			if content != testcontent {
				t.Fatalf("For %s:\nExpected:\n`%s`\nGot:\n`%s`\n", name, testcontent, content)
			}

			continue
		}

		// Not a web/* test yarn.
		content, ok := web.Get(name)

		if ok {
			t.Fatalf("Found unexpected key: %s in Sub(web): %s", name, content)
		}
	}
}

func TestList(t *testing.T) {

	files := testdata.List()

	for _, file := range files {
		_, ok := testyarns[file]
		if !ok {
			t.Fatalf("Unexpect file in list: `%s`", file)
		}
	}

	for _, key := range goodkeys {

		for _, file := range files {
			if file == key {
				goto next
			}
		}

		t.Fatalf("Missing %s in list", key)
	next:
		continue
	}

}

func TestListSub(t *testing.T) {

	testprefix := "web"

	web := testdata.Sub(testprefix)

	files := web.List()

	// check the return list
	for _, file := range files {

		if !strings.HasPrefix(file, testprefix) {
			t.Fatalf("Unexpected file %s in sub: %s", file, testprefix)
		}

		_, ok := testyarns[file]
		if !ok {
			t.Fatalf("Unexpect file in list: `%s`", file)
		}
	}

	// check if the list is complete

	for _, key := range goodkeys {

		// we don't expect things outside of the prefix
		// we have already test that there isn't any file that doesn't
		// match the prefix in the previous step, so we just continue here.
		if !strings.HasPrefix(key, testprefix) {
			continue
		}

		for _, file := range files {
			if file == key {
				goto next
			}
		}

		t.Fatalf("Missing %s in list", key)
	next:
		continue
	}

}

func TestWalk(t *testing.T) {

	files := map[string]string{}

	testdata.Walk("**", func(path, content string) {
		files[path] = content
	})

	for name, testcontent := range testyarns {
		content, ok := files[name]
		if !ok {
			t.Fatalf("Missing %s", name)
		}
		if content != testcontent {
			t.Fatalf("For %s:\nExpected:\n`%s`\nGot:\n`%s`\n", name, testcontent, content)
		}
	}

	for path := range files {

		found := false

		for name := range testyarns {
			if path == name {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("unexpected file %s.", path)
		}

	}

}

func TestWalkSub(t *testing.T) {
	// TODO
}

func TestAll(t *testing.T) {

	files := testdata.All()

	for name, testcontent := range testyarns {
		content, ok := files[name]
		if !ok {
			t.Fatalf("Missing %s", name)
		}
		if content != testcontent {
			t.Fatalf("For %s:\nExpected:\n`%s`\nGot:\n`%s`\n", name, testcontent, content)
		}
	}

	for path := range files {

		found := false

		for name := range testyarns {
			if path == name {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("unexpected file %s.", path)
		}

	}

}

func TestAllSub(t *testing.T) {
	// TODO
}
