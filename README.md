# GoWiki
Go implementation of a simple Wiki.

This code is based on one of the official Golang Tutorials here: https://golang.org/doc/articles/wiki/

## Covered in this tutorial :) For fun to make conflict

- Creating a data structure with load and save methods
- Using the net/http package to build web applications
- Using the html/template package to process HTML templates
- Using the regexp package to validate user input
- Using closures


### Wiki page data structure

```
type Page struct {
    Title string
    Body  []byte
}
```

### Building and running the application

```
$ go build wiki.go
$ ./wiki
```
