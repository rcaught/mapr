# mapr [![Build Status](https://travis-ci.org/rcaught/mapr.svg?branch=master)](https://travis-ci.org/rcaught/mapr) [![GoDoc](https://godoc.org/github.com/rcaught/mapr?status.svg)](https://godoc.org/github.com/rcaught/mapr)

Apply commands to values of structured data, with optional key filtering.

## Installation
##### Go
```
$ go get github.com/rcaught/mapr/...
```
##### MacOS
```
$ curl -Ls https://github.com/rcaught/mapr/releases/latest/download/macos.zip > /tmp/mapr.zip
$ unzip /tmp/mapr.zip -d /usr/local/bin
```
##### Linux
```
$ curl -Ls https://github.com/rcaught/mapr/releases/latest/download/linux.zip > /tmp/mapr.zip
$ unzip /tmp/mapr.zip -d /usr/local/bin
```

## Usage
```
$ # Uppercase all values
$ cat fixtures/data.json | mapr "echo {{value}} | tr a-z A-Z"

$ # Find all keys with a suffix of _enc.  For each of these key/values, remove the key suffix and decrypt the value in the resulting JSON
$ cat fixtures/data.json | mapr --key-filter-type=suffix --key-filter=_enc --key-filter-strip "awscli kms decrypt --ciphertext-blob {{value}}"
```
