# perfeng-search

## Setup

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```shell script
$ export PROJECT_PATH="$(go env GOROOT)/src/github.com/usfcs-perf-eng-s20/"
$ mkdir -p "$PROJECT_PATH"
$ git clone git@github.com:usfcs-perf-eng-s20/heroku-project-search.git $PROJECT_PATH
$ cd $PROJECT_PATH
```

## Obtain DB credentials
```shell script
$ heroku config -a perfeng-go-search -s >> .env
```

Ask Daniel, Zini or Olivia on Slack for Heroku app access.

## Running Locally

```shell script
$ make local
```

or alternatively:

```shell script
$ go build -o bin/project-go-search -v .
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
