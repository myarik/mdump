# Database backup manager application

**mdump** is a simple database backup manager. 

## Installation

### Docker container

The simple way to install an application is using Docker

```shell
docker build -t mdump:latest .
```

### From Source

mdump is written in the Go programming language and you need at least Go version 1.16. Mdump may also work with older versions of Go, but that’s not supported. 
In order to build mdump from source, execute the following steps:

```shell
$ git clone https://github.com/myarik/mdump.git
[...]

$ cd mdump

$ go build -o mdump -ldflags "-s -w" ./main.go
```

The binary requires database dump tools (pg_dump, etc...)

## Backing up

### Local 

Backing up databases to the local storage

```shell
$ /usr/bin/mdump pgdump local --pg_uri postgresql://[user[:password]@][netloc][:port] --path /backup
```

### Amazon S3

Backing up databases to an Amazon S3 bucket. You must first set up the following environment variables with the credentials.

```shell
$ export AWS_ACCESS_KEY_ID=<MY_ACCESS_KEY>
$ export AWS_SECRET_ACCESS_KEY=<MY_SECRET_ACCESS_KEY>
$ export AWS_REGION=eu-west-2 
```

Backing up databases to the S3

```shell
$ /usr/bin/mdump pgdump s3 --pg_uri postgresql://[user[:password]@][netloc][:port] --aws-bucket bucket_name --aws-bucket-key bucket_key
```

## Supported databases

- PostgreSQL