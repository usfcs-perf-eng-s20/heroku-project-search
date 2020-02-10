# perfeng-search

## Setup

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```shell script
$ export PROJECT_PATH="$(go env GOROOT)/src/github.com/usfcs-perf-eng-s20/"
$ git clone git@github.com:usfcs-perf-eng-s20/project-go-search.git $PROJECT_PATH
$ cd $PROJECT_PATH
```

## Running Locally

```shell script
$ make local
```

or alternatively:

```shell script
$ go build -o bin/go-search -v .
github.com/mattn/go-colorable
gopkg.in/bluesuncorp/validator.v5
golang.org/x/net/context
github.com/heroku/x/hmetrics
github.com/gin-gonic/gin/render
github.com/manucorporat/sse
github.com/heroku/x/hmetrics/onload
github.com/gin-gonic/gin/binding
github.com/gin-gonic/gin
github.com/heroku/project-go-search
$ heroku local
```

The app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```shell script
$ git push heroku master
$ heroku open
```


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
