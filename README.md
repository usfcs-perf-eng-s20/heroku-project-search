# perfeng-search

## Setup

Make sure you have [Go](http://golang.org/doc/install) version 1.13 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```shell script
$ export PROJECT_PATH="$(go env GOROOT)/src/github.com/usfcs-perf-eng-s20/"
$ mkdir -p "$PROJECT_PATH"
$ git clone git@github.com:usfcs-perf-eng-s20/heroku-project-search.git $PROJECT_PATH
$ cd $PROJECT_PATH
```

## Obtain DB credentials
```shell script
$ heroku config -a perfeng-search -s >> .env
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
$ git push origin master
$ heroku open -a perfeng-go-search
```


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)


## Database
Login to remote Postgres instance
```shell script
psql -h <host> -p <port> -W -U <user> -d <db name>
```

Create schema:
```sql
CREATE TABLE dvds (
id bigint PRIMARY KEY,
title VARCHAR (300) NOT NULL,
studio VARCHAR(50) NOT NULL,
price VARCHAR(10) NOT NULL,
rating VARCHAR(30) NOT NULL,
year VARCHAR(10) NOT NULL,
genre VARCHAR(30) NOT NULL,
upc VARCHAR(30) NOT NULL);
```

Create a CSV with the identical schema
```text
id,title,studio,price,rating,year,genre,upc
316270,Innocent Man (1989/ Kino Lorber Studio Classics/ Special Edition),Kino Lorber Studio Classics,$14.95,R,1989,Drama,738329235635
316271,Innocent Man (1989/ Kino Lorber Studio Classics/ Special Edition/ Blu-ray),Kino Lorber Studio Classics,$19.95,R,1989,Drama,738329235642
...
```

Populate table from a csv file
```sql
\copy dvds FROM <path to file> with (format csv,header true, delimiter ',');
```