package bitbucket

import (
	"reflect"
	"testing"
)

func TestMatchPipeline(t *testing.T) {
	config, err := ParseString(pipelineYaml)
	if err != nil {
		t.Error(err)
		return
	}

	got, want := config.Pipeline("refs/tags/release-1.0", "master"), config.Pipelines.Tags["release-*"]
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expect release-* pipeline matches release-1.0")
	}

	got, want = config.Pipeline("", "staging"), config.Pipelines.Branches["staging"]
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expect staging pipeline matches staging branch")
	}

	got, want = config.Pipeline("refs/tags/v1.0.0", "master"), config.Pipelines.Default
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expect default branch used when no match is found")
	}
}

var pipelineYaml = `
image: node:latest

pipelines:
  default:
    - step:
        script:
          - npm install
          - npm test
  tags:
    release-*:
      - step:
          script:
            - npm install
            - npm test
            - npm run release
  branches:
    staging:
      - step:
          script:
            - echo "Clone all the things!"
`
