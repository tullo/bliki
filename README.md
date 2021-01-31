# Bliki

A simple, self-containig blog engine that consists of just one binary. The binary contains html templates and an embedded sqlite datebase, thus no external dependencies are required.

## Build

```text
cd && GO111MODULE=on go get github.com/gobuffalo/packr/v2/packr2@v2.8.1
make build
``` 

## Usage

```text
./bin/bliki
``` 

The engine expects the database file `bliki.db` in the current directory.
The file is created automatically on the engines first start.

## URLs

```text
Index | http://{HOST}:3000
Admin | http://{HOST}:3000/admin
```

## Environment variables

Add the following evironment variables to an `.env` file in current directory:

```text
export USERNAME=super    # username for admin login
export PASSWORD=s3cr3t   # password for admin login
export PORT=3000         # http server port
```
