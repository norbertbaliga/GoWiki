# GoWiki
Go implementation of a simple Wiki.

This code is based on one of the official Golang Tutorials here: https://golang.org/doc/articles/wiki/

Added features:

- Index page at `/` path to show list of available pages.
- Serve static content under `/js` and `/css`.
- Create/Update page from the Index page.
- Delete function. It is possible to delete pages from the index page and the view page.
- Search function. Filter Wiki pages based on the searchterm provided.

Running the application will launch a webserver listening on the port 8888 by default. The default port can be changed by setting a `GOWIKI_LISTEN_PORT` environment variable. 

After running the application open the `http://localhost:8888/` URL in a browser or any HTTP client to access the Index page of GoWiki site. Here there is a list of available pages that can be viewed and it is possible to create a new page. Navigation between the Index page and View/Edit are possible via links placed on the sites. Full CRUD functionalites have been implemented.

Directly visiting the `http://localhost:8888/view/<page title>` URL it is possible to view the content of the `<page title>.txt` text file under the pages folder defined by `GOWIKI_PAGES_PATH` environment variable (`./pages` by default if it is not set). If the title does not exist a `Page Not Found` page will be presented with a button to create a page with the given title. Similar to the View directly visiting the `http://localhost:8888/edit/<page title>` URL it is possible to edit the content of the title if exists, otherwise new page can be created.


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

1. List pages

At the root path '/' the Index page shows the available Wiki pages. Lists the .txt files under pages folder defined by `GOWIKI_PAGES_PATH` environment variable (`./pages` by default if it is not set).

2. Create pages

By visiting the `/view/<page title>` or `/edit/<page title>` URL with a page title that does not exist a new page can be created. In addition to that it is possible to provide a title on the Index page and clicking on the 'Create new page' button.

3. View pages

To show the content of a page visit `/view/<page title>` URL. However, on the Index page the titles are links that brings the user to this view. 

4. Edit pages

To edit a Title visit the `/edit/<page title>` URL or use the 'edit' link on a page viewed.

5. Delete pages

Delete links are available on the Index page next to the titles and on the view page. This removes the given `.txt` file under pages folder defined by `GOWIKI_PAGES_PATH` environment variable (`./pages` by default if it is not set).

6. Search pages

It is possible to filter the listed Wiki pages on the Index site by providing a searchterm in the `Search pages` field or adding `?q=<searchterm>` GET parameter to the URL directly.


### Building and running the application

```
$ go build wiki.go
$ ./wiki
```
This will run the GoWiki application listening on the default port 8888 and will store the pages created under the `./pages` folder. This can be changed by setting the following environment variables.

```
$ export GOWIKI_LISTEN_PORT=<specific port number>
$ export GOWIKI_PAGES_PATH=<path to folder>
$ ./wiki
```

### Building Docker image and running the application in container

```
$ docker build -t <gowiki:tag> .
$ docker run -it --rm -p 8888:8888 <gowiki:tag>
```

Choose an appropriate tag `<gowiki:tag>` for the image based on the application version.

By using the `pages_folder` build argument it is possible to change the default path where the wiki pages will be stored in the container filesystem. For example, to view and store pages from the `/home/pages` build the image as follows:

```
$ docker build -t <gowiki:tag> --build-arg pages_folder=/home/pages .
```

Running the container with modified port and pages folder path specify the environment variable as follows:

```
$ docker run -it --rm -p 8888:<my port> -e GOWIKI_LISTEN_PORT=<my port> -e GOWIKI_PAGES_PATH=<my folder path> <gowiki:tag>
```