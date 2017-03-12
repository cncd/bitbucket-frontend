package bitbucket

import "testing"

func TestParse(t *testing.T) {
	config, err := ParseString(sample)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := config.Clone.Depth, 25; want != got {
		t.Errorf("Wanted clone depth %d, got %d", want, got)
	}

	if want, got := len(config.Pipelines.Default.Steps), 2; want != got {
		t.Errorf("Wanted default.step length %d, got %d", want, got)
		t.FailNow()
	}

	if want, got := config.Pipelines.Default.Steps[0].Image, "golang:1.7"; want != got {
		t.Errorf("Wanted default.step.0.image %s, got %s", want, got)
	}

	if want, got := len(config.Pipelines.Default.Steps[0].Script), 2; want != got {
		t.Errorf("Wanted default.step.0.step length %d, got %d", want, got)
	}

	if want, got := config.Pipelines.Default.Steps[0].Script[0], "go build"; want != got {
		t.Errorf("Wanted default.step.0.step.0 equal to %s, got %s", want, got)
	}

	if want, got := config.Pipelines.Default.Steps[1].Image, ""; want != got {
		t.Errorf("Wanted default.step.0.image equal to %q, got %s", want, got)
	}
}

var sample = `
image: node:latest

clone:
  depth: 25

pipelines:
  default:
    - step:
        image: golang:1.7
        script:
          - go build
          - go test
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
