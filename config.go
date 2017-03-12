package bitbucket

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

	//
	// the below metadata structs are copied from github.com/cncd/config
	// please do not change.
	//

	// MetaData defines runtime metadata.
	MetaData struct {
		Repo Repo
		Curr Build
		Prev Build
		Job  Job
		Sys  System
	}

	// Repo defines runtime metadata for a repository.
	Repo struct {
		Name    string
		Link    string
		Remote  string
		Private bool
	}

	// Build defines runtime metadata for a build.
	Build struct {
		Number   int
		Created  int64
		Started  int64
		Finished int64
		Status   string
		Event    string
		Link     string
		Target   string

		Commit Commit
	}

	// Commit defines runtime metadata for a commit.
	Commit struct {
		Sha     string
		Ref     string
		Refspec string
		Branch  string
		Message string
		Author  Author
	}

	// Author defines runtime metadata for a commit author.
	Author struct {
		Name   string
		Email  string
		Avatar string
	}

	// Job defines runtime metadata for a job.
	Job struct {
		Number int
		Matrix map[string]string
	}

	// System defines runtime metadata for a ci/cd system.
	System struct {
		Name string
		Link string
		Arch string
	}
)

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
