package yarn

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

var (
	testyarns map[string]string
	goodkeys  []string
	badkeys   []string
	sqls      Yarn
)

func TestMain(m *testing.M) {
	testyarns = map[string]string{
		"insert.sql":    "INSERT INTO users (id, name, email)\nVALUES ($1, $2, $3)\n",
		"query_all.sql": "SELECT\nid,\nname,\nFROM users\n",
	}

	goodkeys = []string{"insert.sql", "query_all.sql"}

	badkeys = []string{"something.json", "random", "sql", "insert", "nope", "none"}
	sqls = Must(http.Dir("testdata"), "*.sql")

	os.Exit(m.Run())
}

func TestMustHave(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("Didn't panic for bad keys: %s", badkeys)
		}
		if r == fmt.Sprintf(missingYarn, badkeys) {
			return
		}
		panic(r)
	}()
	sqls.MustHave(badkeys...)
	//TODO: Handle this properly.
	sqls.MustHave(goodkeys...)
}

func TestHas(t *testing.T) {

	err := sqls.Has(goodkeys...)
	if err != nil {
		t.Fatal(err)
	}

	if err := sqls.Has(badkeys...); err == nil {
		t.Fatal("Expected error. Got nothing.")
	}
}

func TestMust(t *testing.T) {
	for name, testcontent := range testyarns {
		func(name, testcontent string) {
			defer func() {
				r := recover()
				if r != nil {
					if r == fmt.Sprintf(missingYarn, name) {
						t.Fatalf("Missing using MUST %s", name)
						return
					}
					panic(r)
				}
			}()
			content := sqls.Must(name)
			if content != testcontent {
				t.Fatalf("For %s:\nExpected:`%s`\nGot:`%s`\n", name, testcontent, content)
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
				if r != fmt.Sprintf(missingYarn, name) {
					panic(r)
				}
			}()
			sqls.Must(name)

		}(name)
	}
}

func TestGet(t *testing.T) {
	for name, testcontent := range testyarns {
		content, ok := sqls.Get(name)
		if !ok {
			t.Fatalf("Missing %s", name)
		}
		if content != testcontent {
			t.Fatalf("For %s:\nExpected:`%s`\nGot:`%s`\n", name, testcontent, content)
		}
	}

	for _, name := range badkeys {
		if _, ok := sqls.Get(name); ok {
			t.Fatalf("Got OK for unexpected `%s` key.", name)
		}
	}
}
