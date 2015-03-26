package yarn

import (
	"fmt"
	"net/http"
	"testing"
)

func TestYarn(t *testing.T) {
	testsqls := map[string]string{
		"insert.sql":    "INSERT INTO users (id, name, email)\nVALUES ($1, $2, $3)\n",
		"query_all.sql": "SELECT\nid,\nname,\nFROM users\n",
	}
	sqls := Must(http.Dir("testdata"), "*.sql")

	err := sqls.Has("insert.sql", "query_all.sql")
	if err != nil {
		t.Fatal(err)
	}

	if err := sqls.Has("nope", "not this one"); err == nil {
		t.Fatal("Expected error. Got nothing.")
	}

	for name, testcontent := range testsqls {
		content, ok := sqls.Get(name)
		if !ok {
			t.Fatalf("Missing %s", name)
		}
		if content != testcontent {
			t.Fatalf("For %s:\nExpected:`%s`\nGot:`%s`\n", name, testcontent, content)
		}

		defer func() {
			r := recover()
			if r != nil && r != fmt.Sprintf(missingYarn, name) {
				panic(r)
			}
		}()
		content = sqls.Must(name)
		if content != testcontent {
			t.Fatalf("For %s:\nExpected:`%s`\nGot:`%s`\n", name, testcontent, content)
		}
	}

	for _, name := range []string{"something.json", "random", "sql", "insert", "nope", "none"} {
		if _, ok := sqls.Get(name); ok {
			t.Fatal("Got OK for unexpected `%s` key.", name)
		}

		func(name string) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("Must didn't panic for unexpected `%s` key.", name)
				}
				if r != fmt.Sprintf(missingYarn, name) {
					panic(r)
				}
			}()
			sqls.Must(name)

		}(name)
	}
}
