# goredirects [![Report card](https://goreportcard.com/badge/bramp.net/goredirects)](https://goreportcard.com/report/bramp.net/goredirects) [![GoDoc](https://godoc.org/bramp.net/goredirects?status.svg)](https://godoc.org/bramp.net/goredirects)

by Andrew Brampton ([bramp.net](https://bramp.net))

goredirects enables the use of a vanity redirect domain in your go package
imports. For example, instead of using `import "github.com/example/package"` you
could use a vanity domain, and `import "example.com/package"`, yet still host the
source code on GitHub.

Specifically, this tool creates a set of HTML files containing the go-imports
meta tags, for each of your projects. This uses the [remote import
paths](https://golang.org/cmd/go/#hdr-Remote_import_paths) feature of `go get`
command to redirect from your vanity domain to GitHub.com.

## Example
To create a set of static HTML redirects:

```bash
$ go install bramp.net/goredirects@latest

$ goredirects
Usage: goredirects <domain> <output dir>

$ goredirects bramp.net outputdir
# Looking under $GOROOT/src/bramp.net for all packages
...
```

This will search your $GOROOT/src/<domain> for all packages, and create static
HTML into the outputdir for each package.

To read more about how this tool works, checkout my [blog article](https://blog.bramp.net/post/2017/10/02/vanity-go-import-paths/) on the topic.

## Developing

Before committing, please run:

```shell
go fmt

go test

go vet ./...

go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...

go install golang.org/x/lint/golint@latest
golint ./...
```

### Updating test data

There are test git repos under `test/input`. Git forbides a `.git` directory to be checked in, so we put them all in a tar file, that is extracted during testing.

Extract the tar

```shell
tar -xvf test/input.tar.gz test/input
```

Create a new tar file:

```shell
find test/input -print0 | LC_ALL=C sort -z |
tar --no-recursion --null -T - \
    --no-xattrs \
    --no-recursion \
    --options '!timestamp' \
    -cvzf test/input.tar.gz
```

## Licence (Apache 2)

*This is not an official Google product (experimental or otherwise), it is just
code that happens to be owned by Google.*

```
Copyright 2017 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
