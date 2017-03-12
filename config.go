package bitbucket

import (
	"path"
	"strings"
)

// see https://confluence.atlassian.com/bitbucket/configure-bitbucket-pipelines-yml-792298910.html#Configurebitbucket-pipelines.yml-ci_branches

type (
	// Config defines the pipeline configuration.
	Config struct {
		// Image specifies the Docker image with
		// which we run your builds.
		Image string

		// Clone defines the depth of Git clones
		// for all pipelines.
		Clone struct {
			Depth int
		}

		// Pipeline defines the pipeline configuration
		// which includes a list of all steps for default,
		// tag, and branch-specific execution.
		Pipelines struct {
			Default  Stage
			Tags     map[string]Stage
			Branches map[string]Stage
		}
	}

	// Stage contains a list of steps executed
	// for a specific branch or tag.
	Stage struct {
		Name  string
		Steps []*Step
	}

	// Step defines a build execution unit.
	Step struct {
		// Image specifies the Docker image with
		// which we run your builds.
		Image string

		// Script contains the list of bash commands
		// that are executed in sequence.
		Script []string
	}
)

// Pipeline returns the pipeline stage that best matches the branch
// and ref. If there is no matching pipeline specific to the branch
// or tag, the default pipeline is returned.
func (c *Config) Pipeline(ref, branch string) Stage {
	// match pipeline by tag name
	tag := strings.TrimPrefix(ref, "refs/tags/")
	for pattern, pipeline := range c.Pipelines.Tags {
		if ok, _ := path.Match(pattern, tag); ok {
			return pipeline
		}
	}
	// match pipeline by branch name
	for pattern, pipeline := range c.Pipelines.Branches {
		if ok, _ := path.Match(pattern, branch); ok {
			return pipeline
		}
	}
	// use default
	return c.Pipelines.Default
}

// UnmarshalYAML implements custom parsing for the stage section of the yaml
// to cleanup the structure a bit.
func (s *Stage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	in := []struct {
		Step *Step
	}{}
	err := unmarshal(&in)
	if err != nil {
		return err
	}
	for _, step := range in {
		s.Steps = append(s.Steps, step.Step)
	}
	return nil
}
