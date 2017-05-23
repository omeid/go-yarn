# Catalog
[![GoDoc](https://godoc.org/github.com/omeid/go-yarn?status.svg)](https://godoc.org/github.com/omeid/go-yarn/catalog) [![Coverage Status](https://coveralls.io/repos/omeid/go-yarn/badge.png)](https://coveralls.io/r/omeid/go-yarn) [![Build Status](https://drone.io/github.com/omeid/go-yarn/status.png)](https://drone.io/github.com/omeid/go-yarn/latest)

A simple multi yarn per file parser. It allows you to use single files as a catalog of Yarns.
Some usecases include sql migration, sql statements, redis queries with lua, et al.


## Syntax

A Catalog file is a list of content encapsulated in Start and End separator lines.  

Catalog supports and automatically detects four different syntax for separator lines, each of which will constitue a comment line in one or many programming languages, this allows you to use the appropriate filetype and extension and leverage your text editor tooling (linting, syntax highlighting, et al).

The supported comment types are lines that start with the C-like (`//`), Bash-like (`##`), and Ada-like (`--`) tokens, or a line starting with `/*` and ending with `*/` (C-block style).

The following files are all valid catalogs:


```sql

-- start: select_all_users
SELECT * from users
-- end: select_all_users

-- start: pages
SELECT * from pages
-- end: pages
```

```lua
-- start: HSETEX
local res = redis.call("HSET", KEYS[1], ARGV[2], ARGV[3])
redis.call("EXPIRE", KEYS[1], ARGV[1])
-- end: HSETEX 

```
```bash
## start: scriptB

set -xe
function Scripting() {
  # Do lots of unholy things here.
}

# Run the script
Scripting

## end: scriptB


# More scripting
echo "Why do you run bash from Go? Double You Tee If?"
## end: scriptB
```




## Usage

The simplest way to use this package is with yarn, you simply setup a yarn that includes your catalog-file, then Catalog provides you with a Yarn instance that includes all the files in your catalog.


```go

// Assuming the first example is in a file called catalog.sql in the
// package loaded by the the yarn instance called scripts.
sqls, err := catalog.Parse(scripts, "catalog.sql")
// Deal with the error

//Now sqls will have every file from your catalog. Simply load them
err := sqls.Have("select_all_users", "get_all_pages")
// Deal with the error

// Here you can just use the standard yarn api to get the content of each file in your catalog.
res, err := sql.Exec(sqls.Must("select_all_users"))
// And so forth.
```


For more details, see the [API Docs](https://godoc.org/github.com/omeid/go-yarn/catalog).
