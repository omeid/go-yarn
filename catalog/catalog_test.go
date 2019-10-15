package catalog

import "testing"

var goodCases = []struct {
	sources []string
	yarns   map[string]string
}{
	{
		sources: []string{`
-- start: select_all_users
SELECT * from users
-- end: select_all_users

-- start: pages hello
SELECT * from pages
-- end: pages hello

`,
			`
## start: select_all_users
SELECT * from users
## end: select_all_users

## start: pages hello
SELECT * from pages
## end: pages hello

`,
			`
// start: select_all_users
SELECT * from users
// end: select_all_users

// start: pages hello
SELECT * from pages
// end: pages hello

`,
			`
/* start: select_all_users */
SELECT * from users
/* end: select_all_users */

/* start: pages hello */
SELECT * from pages
/* end: pages hello */

`,
		},
		yarns: map[string]string{
			"select_all_users": "SELECT * from users",
			"pages hello":      "SELECT * from pages",
		},
	},
}

func TestCatalog(t *testing.T) {

	for _, c := range goodCases {

		for _, source := range c.sources {
			y, err := ParseString(source)

			if err != nil {
				t.Fatal(err)
			}

			for name, expect := range c.yarns {
				content, found := y.Get(name)
				if !found {
					t.Fatalf("Expected %s but was not found.", name)
				}

				if content != expect {
					t.Fatalf("For: %v\nExpected:\n%v\nGot:\n%v", name, expect, content)
				}
			}
		}
	}
}

var testfiles = map[string]map[string]string{
	"testdata/sample.sql": map[string]string{
		"select_all_users": "SELECT * from users;",
		"pages hello":      "SELECT * from pages;",
	},
}

func TestCatalogFile(t *testing.T) {

	for filename, expect := range testfiles {
		y, err := ParseFile(filename)

		if err != nil {
			t.Fatalf("Unexpected err: %v", err)
		}

		for name, expect := range expect {
			content, found := y.Get(name)
			if !found {
				t.Fatalf("Expected %s but was not found.", name)
			}

			if content != expect {
				t.Fatalf("For: %v\nExpected:\n%v\nGot:\n%v", name, expect, content)
			}
		}
	}
}
