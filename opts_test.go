package bitbucket

import (
	"testing"

	"github.com/cncd/pipeline/pipeline/frontend"
)

func TestWithVolumes(t *testing.T) {
	compiler := NewCompiler(
		WithVolumes(
			"/tmp:/tmp",
			"/foo:/foo",
		),
	)
	if compiler.volumes[0] != "/tmp:/tmp" || compiler.volumes[1] != "/foo:/foo" {
		t.Errorf("TestWithVolumes must set default volumes")
	}
}

func TestWithWorkspace(t *testing.T) {
	compiler := NewCompiler(
		WithWorkspace(
			"/pipeline",
			"src/github.com/octocat/hello-world",
		),
	)
	if compiler.base != "/pipeline" {
		t.Errorf("WithWorkspace must set the base directory")
	}
	if compiler.path != "src/github.com/octocat/hello-world" {
		t.Errorf("WithWorkspace must set the path directory")
	}
}

func TestWithPrefix(t *testing.T) {
	if NewCompiler(WithPrefix("drone_")).prefix != "drone_" {
		t.Errorf("WithPrefix must set the prefix")
	}
}

func TestWithMetadata(t *testing.T) {
	metadata := frontend.Metadata{
		Repo: frontend.Repo{
			Name:    "octocat/hello-world",
			Private: true,
			Link:    "https://github.com/octocat/hello-world",
			Remote:  "https://github.com/octocat/hello-world.git",
		},
	}
	compiler := NewCompiler(
		WithMetadata(metadata),
	)
	if compiler.env["CI_REPO_NAME"] != metadata.Repo.Name {
		t.Errorf("WithMetadata must set CI_REPO_NAME")
	}
	if compiler.env["CI_REPO_LINK"] != metadata.Repo.Link {
		t.Errorf("WithMetadata must set CI_REPO_LINK")
	}
	if compiler.env["CI_REPO_REMOTE"] != metadata.Repo.Remote {
		t.Errorf("WithMetadata must set CI_REPO_REMOTE")
	}
}

func TestWithLocal(t *testing.T) {
	if NewCompiler(WithLocal(true)).local == false {
		t.Errorf("WithLocal true must enable the local flag")
	}
	if NewCompiler(WithLocal(false)).local == true {
		t.Errorf("WithLocal false must disable the local flag")
	}
}

func TestWithNetrc(t *testing.T) {
	compiler := NewCompiler(
		WithNetrc(
			"octocat",
			"password",
			"github.com",
		),
	)
	if compiler.env["CI_NETRC_USERNAME"] != "octocat" {
		t.Errorf("WithNetrc should set CI_NETRC_USERNAME")
	}
	if compiler.env["CI_NETRC_PASSWORD"] != "password" {
		t.Errorf("WithNetrc should set CI_NETRC_PASSWORD")
	}
	if compiler.env["CI_NETRC_MACHINE"] != "github.com" {
		t.Errorf("WithNetrc should set CI_NETRC_MACHINE")
	}
}

func TestWithProxy(t *testing.T) {
	// alter the default values
	noProxy := "foo.com"
	httpProxy := "bar.com"
	httpsProxy := "baz.com"

	testdata := map[string]string{
		"no_proxy":    noProxy,
		"NO_PROXY":    noProxy,
		"http_proxy":  httpProxy,
		"HTTP_PROXY":  httpProxy,
		"https_proxy": httpsProxy,
		"HTTPS_PROXY": httpsProxy,
	}
	compiler := NewCompiler(
		WithProxy(httpProxy, httpsProxy, noProxy),
	)
	for key, value := range testdata {
		if compiler.env[key] != value {
			t.Errorf("WithProxy should set %s=%s", key, value)
		}
	}
}

func TestWithEnviron(t *testing.T) {
	compiler := NewCompiler(
		WithEnviron(
			map[string]string{
				"RACK_ENV": "development",
				"SHOW":     "true",
			},
		),
	)
	if compiler.env["RACK_ENV"] != "development" {
		t.Errorf("WithEnviron should set RACK_ENV")
	}
	if compiler.env["SHOW"] != "true" {
		t.Errorf("WithEnviron should set SHOW")
	}
}
