// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"os"
	"os/exec"
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
			t.Errorf("gitSSHtoHTTPS(%q) = %q, want %q", test.input, got, test.want)
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
			t.Errorf("gitHTTPStoWeb(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestCreate(t *testing.T) {
	redirect := redirectCreator{
		vanity:    "example.com",
		input:     "test/input",
		output:    "test/got",
		gitRemote: "origin",
	}

	// Cleanup once we are done
	//defer os.RemoveAll("test/got")

	if err := redirect.Create(); err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Check if the output files match what was expected
	cmd := exec.Command("diff", "-ur", "test/got", "test/want")
	if out, err := cmd.Output(); err != nil {
		t.Fatalf("diff: %v\n%s", err, out)
	}
}

func extractTestData() error {
	// We use -xvf to extract, and try to ensure .git directories are included.
	// Some tar versions (like on macOS) might have different defaults.
	return exec.Command("tar", "-xzf", "test/input.tar.gz").Run()
}

func TestMain(m *testing.M) {
	if err := extractTestData(); err != nil {
		fmt.Printf("Failed to extract test data: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
