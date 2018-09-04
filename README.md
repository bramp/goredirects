# goredirects [![Build Status](https://img.shields.io/travis/bramp/goredirects.svg)](https://travis-ci.org/bramp/goredirects) [![Coverage](https://img.shields.io/coveralls/bramp/goredirects.svg)](https://coveralls.io/github/bramp/goredirects) [![Report card](https://goreportcard.com/badge/bramp.net/goredirects)](https://goreportcard.com/report/bramp.net/goredirects) [![GoDoc](https://godoc.org/bramp.net/goredirects?status.svg)](https://godoc.org/bramp.net/goredirects)

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
$ go install bramp.net/goredirects

$ goredirects
Usage: goredirects <domain> <output dir>
  -git-remote string
        Git remote name (default "origin")
  -include-vendor
        Include vendor directory

$ goredirects bramp.net outputdir
# Looking under $GOROOT/src/bramp.net for all packages
...
```

This will search your $GOROOT/src/<domain> for all packages, and create static
HTML into the outputdir for each package.

To read more about how this tool works, checkout my [blog article](https://blog.bramp.net/post/2017/10/02/vanity-go-import-paths/) on the topic.

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