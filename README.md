# GoWiki
Go implementation of a simple Wiki.

This code is based on one of the official Golang Tutorials here: https://golang.org/doc/articles/wiki/

Running the application will launch a webserver listening on the port 8080. Visit the `http://localhost:8080/view/<page title>` in a browser or any HTTP client that will show the content of the `<page title>.txt` text file under the `/pages` folder if it already exists. In this case in addition to the content there is also an 'edit' link on the page that can be used to update the content of this file. If the title does not exist a 'Page Not Found' page will be presented with a button to create a page with the given title.


## Covered in this tutorial

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

### Application functions

1. Viewing pages

Visiting the `/view/<page title>` URL will show the page content if the page already exists.

3. Editing pages

TODO

5. Saving pages

TODO


### Building and running the application

```
$ go build wiki.go
$ ./wiki
```
