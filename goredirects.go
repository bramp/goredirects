// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// goredirects creates static HTML that redirects go packages hosted
// on a vanity domain to GitHub.
package main // import "bramp.net/goredirects"

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"net/url"

	git "github.com/go-git/go-git/v5"
)

const html = `<html>
<head>
<meta name="go-import" content="{{.Name}} git {{.RepoURL}}">
<meta http-equiv="refresh" content="0; url={{.SiteURL}}" />
<link rel="canonical" href="{{.SiteURL}}" />
<script>
	window.location.replace({{.SiteURL }});
</script>
</head>
<body>
	<h1>Redirecting to <a href="{{.SiteURL}}">{{.SiteURL}}</a></h1>
</body>
</html>
`

type redirectCreator struct {
	vanity string // The vanity domain
	input string  // The path to the root of all go projects
	output string // The output location
}

// templateData holds the data to be rendered by the template
type templateData struct {
	Name    string // Root name "bramp.net/goredirects"
	RepoURL template.URL // https://github.com/bramp/goredirects.git
	SiteURL template.URL // https://github.com/bramp/goredirects
}

var tmpl = template.Must(template.New("index").Parse(html))

var githubSSHrx = regexp.MustCompile("git@github.com:([^/]*)/(.*).git")
var githubHTTPSrx = regexp.MustCompile("https://github.com/([^/]*)/(.*)(?:.git)?")

// githubSSHtoHTTPS takes a URL to a SSH repo, and returns the equlivant HTTPS url.
// If it is already a valid Github HTTPS URL just return it.
func githubSSHtoHTTPS(url string) string {
	matches := githubSSHrx.FindStringSubmatch(url)
	if len(matches) == 3 {
		return fmt.Sprintf("https://github.com/%s/%s.git", matches[1], matches[2])
	}

	// Perhaps it is already a HTTPS URL?
	if githubHTTPSrx.MatchString(url) {
		return url
	}

	// TODO(bramp) Change this to return an error.
	panic(fmt.Sprintf("not a github repo URL %q", url))
}

// githubHTTPStoWeb takes a URL to a HTTPS repo, and returns the site.
func githubHTTPStoWeb(url string) string {
	if !githubHTTPSrx.MatchString(url) {
		panic(fmt.Sprintf("not a github HTTP URL %q", url))
	}
	return strings.TrimSuffix(url, ".git")
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Printf("Failed to stat %q: %s", path, err)
		return false
	}
	return fi.Mode().IsDir()
}

// Create creates all the redirect HTML pages based on the supplied vanity domain.
func (r *redirectCreator) Create() error {
	repos, err := filepath.Glob(filepath.Join(r.input, "*"))
	if err != nil {
		return fmt.Errorf("failed to read repos: %s", err)
	}

	for _, repopath := range repos {
		// Skip files, and hidden directories.
		if !isDir(repopath) || strings.HasPrefix(filepath.Base(repopath), ".") {
			continue
		}

		if err := r.handleRepo(repopath); err != nil {
			log.Printf("%s", err)
		}
	}

	return nil
}

// Returns the HTTPS/Web URL for this repo.
func (redirectCreator) repoURL(repopath string) (string, error) {
	repo, err := git.PlainOpen(repopath)
	if err != nil {
		return "", fmt.Errorf("failed to open: %s", repopath, err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		// TODO Use repo.Remotes() to find one if origin doesn't exist
		return "", fmt.Errorf("failed to get remote %q: %s", "origin", repopath, err)
	}

	urls := remote.Config().URLs
	if len(urls) != 1 {
		return "", fmt.Errorf("expected only one URL %q: %q", repopath, urls)
	}

	return githubSSHtoHTTPS(urls[0]), nil
}


func (r *redirectCreator) handleRepo(repopath string) error {
	name, err := filepath.Rel(r.input, repopath)
	if err != nil {
		return fmt.Errorf("failed to lookup repo name: %s", err)
	}

	vanityName := url.URL{
		Host: r.vanity,
		Path: name,
	}

	repo, err := r.repoURL(repopath)
	if err != nil {
		return fmt.Errorf("failed to lookup repo url: %s", err)
	}

	// TODO make these use URLs not strings
	repoURL, err := url.Parse(repo);
	if err != nil {
		return fmt.Errorf("failed to parse repo url $q: %s", repo, err)
	}

	site := githubHTTPStoWeb(repo)
	siteURL, err := url.Parse( site );
	if err != nil {
		return fmt.Errorf("failed to parse repo url %q: %s", site, err)
	}

	data := templateData{
		Name:    strings.TrimPrefix(vanityName.String(), "//"),
		RepoURL: template.URL(repoURL.String()),
		SiteURL: template.URL(siteURL.String()),
	}

	r.writeHTML(name, data)

	// Find all sub-packages, and create a similar file
	subpackages := make(map[string]bool)
	subpackages[name] = true

	if err := filepath.Walk(repopath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			dir, err := filepath.Rel(r.input, filepath.Dir(path))
			if err != nil {
				return err
			}
			if _, found := subpackages[dir]; !found {
				r.writeHTML(dir, data)
			}
			subpackages[dir] = true
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk for subpackages %q: %s", repopath, err)
	}

	return nil
}

func (r *redirectCreator) writeHTML(name string, data templateData) error {
	path := filepath.Join(r.output, name)

	log.Printf("Writing %q\n", path)

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("Failed to mkdir %q: %s", path, err)
	}

	f, err := os.Create(filepath.Join(path, "index.html"))
	if err != nil {
		return fmt.Errorf("failed to create %s/index.html: %s", path, err)
	}
	return tmpl.Execute(f, data)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <domain> <input dir> <output dir>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 3 {
		flag.Usage()
		os.Exit(1)
	}

	redirect := redirectCreator{
		vanity: flag.Arg(0),
		input: flag.Arg(1),
		output: flag.Arg(2),
	}
	redirect.Create()
}
