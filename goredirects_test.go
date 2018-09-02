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
//
package main

import (
	"testing"
)

func TestGitSSHtoHTTPS(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"https://github.com/user/project", "https://github.com/user/project"},
		{"https://github.com/user/project.git", "https://github.com/user/project.git"},
		{"git@github.com:user/project.git", "https://github.com/user/project.git"},
		{"https://bitbucket.org/user/project", "https://bitbucket.org/user/project"},
		{"https://bitbucket.org/user/project.git", "https://bitbucket.org/user/project.git"},
		{"git@bitbucket.org:user/project.git", "https://bitbucket.org/user/project.git"},
	}

	for _, test := range tests {
		if got := gitSSHtoHTTPS(test.input); got != test.want {
			t.Errorf("githubSSHtoHTTPS(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestGitHTTPStoWeb(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"https://github.com/user/project", "https://github.com/user/project"},
		{"https://github.com/user/project.git", "https://github.com/user/project"},
		{"https://bitbucket.org/user/project", "https://bitbucket.org/user/project"},
		{"https://bitbucket.org/user/project.git", "https://bitbucket.org/user/project"},
	}

	for _, test := range tests {
		if got := gitHTTPStoWeb(test.input); got != test.want {
			t.Errorf("githubHTTPStoWeb(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
