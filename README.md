# mapr [![Build Status](https://travis-ci.org/rcaught/mapr.svg?branch=master)](https://travis-ci.org/rcaught/mapr)

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

## Examples

Note: The following examples are piped through jq at the end for pretty formatting.

```
$ # Uppercase all values
$ cat fixtures/data.json
{
    "hello": "there",
    "goodbye": {
        "special_enc": "asdf",
        "not_special": "asdf"
    },
    "problem_key_enc": {
        "because_i_have_wrong_value": true
    },
    "special_enc": "asdf"
}
$ cat fixtures/data.json | mapr "echo {{value}} | tr a-z A-Z" | jq .
{
  "goodbye": {
    "not_special": "ASDF",
    "special_enc": "ASDF"
  },
  "hello": "THERE",
  "problem_key_enc": {
    "because_i_have_wrong_value": true
  },
  "special_enc": "ASDF"
}
```

```
$ # Find all keys with a suffix of _enc.  For each of these key/values, remove the key suffix and uppercase the value in the resulting JSON
$ cat fixtures/data.json
{
    "hello": "there",
    "goodbye": {
        "special_enc": "asdf",
        "not_special": "asdf"
    },
    "problem_key_enc": {
        "because_i_have_wrong_value": true
    },
    "special_enc": "asdf"
}

$ cat fixtures/data.json | mapr --key-filter-type=suffix --key-filter=_enc --key-filter-strip "echo {{value}} | tr a-z A-Z" | jq .
{
  "goodbye": {
    "not_special": "asdf",
    "special": "ASDF"
  },
  "hello": "there",
  "problem_key_enc": {
    "because_i_have_wrong_value": true
  },
  "special": "ASDF"
}
```
