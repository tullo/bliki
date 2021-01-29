# Justblog

A simple, self-containig blog engine that consists of just one binary. The binary contains html templates and an embedded sqlite datebase, thus no external dependencies are required. 

## Build

```text
go get github.com/gobuffalo/packr/packr
make
``` 

## Usage

```text
./bin/justblog
``` 

The engine expects the database file `justblog.db` in the current directory. The file is created automatically on the engines first start.

## URLs

```text
Index | http://{HOST}:3000
Admin | http://{HOST}:3000/admin
```

## Environment variables

Add the following evironment variables to an .env file in current directory:

```text
export USERNAME=ralf       # username for admin login
export PASSWORD=password   # password for admin login
export PORT=3000           # http server port
```

