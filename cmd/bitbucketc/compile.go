package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cncd/config-bitbucket"
	"github.com/cncd/pipeline/pipeline/frontend"

	"github.com/urfave/cli"
)

var compileCommand = cli.Command{
	Name:   "compile",
	Usage:  "compile the yaml file",
	Action: compileAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "in",
			Value: "bitbucket-pipelines.yml",
		},
		cli.StringFlag{
			Name:  "out",
			Value: "pipeline.json",
		},
		cli.StringSliceFlag{
			Name: "volumes",
		},
		cli.StringSliceFlag{
			Name: "privileged",
			Value: &cli.StringSlice{
				"plugins/docker",
				"plugins/gcr",
				"plugins/ecr",
			},
		},
		cli.StringFlag{
			Name:  "prefix",
			Value: "pipeline",
		},
		cli.BoolFlag{
			Name: "local",
		},
		//
		// workspace default
		//
		cli.StringFlag{
			Name:  "workspace-base",
			Value: "/workspace",
		},
		cli.StringFlag{
			Name:  "workspace-path",
			Value: "src",
		},
		//
		// netrc parameters
		//
		cli.StringFlag{
			Name:   "netrc-username",
			EnvVar: "CI_NETRC_USERNAME",
		},
		cli.StringFlag{
			Name:   "netrc-password",
			EnvVar: "CI_NETRC_PASSWORD",
		},
		cli.StringFlag{
			Name:   "netrc-machine",
			EnvVar: "CI_NETRC_MACHINE",
		},
		//
		// metadata parameters
		//
		cli.StringFlag{
			Name:   "system-arch",
			Value:  "linux/amd64",
			EnvVar: "CI_SYSTEM_ARCH",
		},
		cli.StringFlag{
			Name:   "system-name",
			Value:  "pipec",
			EnvVar: "CI_SYSTEM_NAME",
		},
		cli.StringFlag{
			Name:   "system-link",
			Value:  "https://github.com/cncd/pipec",
			EnvVar: "CI_SYSTEM_LINK",
		},
		cli.StringFlag{
			Name:   "repo-name",
			EnvVar: "CI_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "repo-link",
			EnvVar: "CI_REPO_LINK",
		},
		cli.StringFlag{
			Name:   "repo-remote-url",
			EnvVar: "CI_REPO_REMOTE",
		},
		cli.StringFlag{
			Name:   "repo-private",
			EnvVar: "CI_REPO_PRIVATE",
		},
		cli.IntFlag{
			Name:   "build-number",
			EnvVar: "CI_BUILD_NUMBER",
		},
		cli.Int64Flag{
			Name:   "build-created",
			EnvVar: "CI_BUILD_CREATED",
		},
		cli.Int64Flag{
			Name:   "build-started",
			EnvVar: "CI_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build-finished",
			EnvVar: "CI_BUILD_FINISHED",
		},
		cli.StringFlag{
			Name:   "build-status",
			EnvVar: "CI_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build-event",
			EnvVar: "CI_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "build-link",
			EnvVar: "CI_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "build-target",
			EnvVar: "CI_BUILD_TARGET",
		},
		cli.StringFlag{
			Name:   "commit-sha",
			EnvVar: "CI_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit-ref",
			EnvVar: "CI_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit-refspec",
			EnvVar: "CI_COMMIT_REFSPEC",
		},
		cli.StringFlag{
			Name:   "commit-branch",
			EnvVar: "CI_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit-message",
			EnvVar: "CI_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "commit-author-name",
			EnvVar: "CI_COMMIT_AUTHOR_NAME",
		},
		cli.StringFlag{
			Name:   "commit-author-avatar",
			EnvVar: "CI_COMMIT_AUTHOR_AVATAR",
		},
		cli.StringFlag{
			Name:   "commit-author-email",
			EnvVar: "CI_COMMIT_AUTHOR_EMAIL",
		},
		cli.IntFlag{
			Name:   "prev-build-number",
			EnvVar: "CI_PREV_BUILD_NUMBER",
		},
		cli.Int64Flag{
			Name:   "prev-build-created",
			EnvVar: "CI_PREV_BUILD_CREATED",
		},
		cli.Int64Flag{
			Name:   "prev-build-started",
			EnvVar: "CI_PREV_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "prev-build-finished",
			EnvVar: "CI_PREV_BUILD_FINISHED",
		},
		cli.StringFlag{
			Name:   "prev-build-status",
			EnvVar: "CI_PREV_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "prev-build-event",
			EnvVar: "CI_PREV_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "prev-build-link",
			EnvVar: "CI_PREV_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "prev-commit-sha",
			EnvVar: "CI_PREV_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "prev-commit-ref",
			EnvVar: "CI_PREV_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "prev-commit-refspec",
			EnvVar: "CI_PREV_COMMIT_REFSPEC",
		},
		cli.StringFlag{
			Name:   "prev-commit-branch",
			EnvVar: "CI_PREV_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "prev-commit-message",
			EnvVar: "CI_PREV_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "prev-commit-author-name",
			EnvVar: "CI_PREV_COMMIT_AUTHOR_NAME",
		},
		cli.StringFlag{
			Name:   "prev-commit-author-avatar",
			EnvVar: "CI_PREV_COMMIT_AUTHOR_AVATAR",
		},
		cli.StringFlag{
			Name:   "prev-commit-author-email",
			EnvVar: "CI_PREV_COMMIT_AUTHOR_EMAIL",
		},
		cli.IntFlag{
			Name:   "job-number",
			EnvVar: "CI_JOB_NUMBER",
		},
		// cli.StringFlag{
		// 	Name:   "job-matrix",
		// 	EnvVar: "CI_JOB_MATRIX",
		// },
	},
}

func compileAction(c *cli.Context) (err error) {
	file := c.Args().First()
	if file == "" {
		file = c.String("in")
	}

	conf, err := bitbucket.ParseFile(file)
	if err != nil {
		return err
	}

	// configure volumes for local execution
	volumes := c.StringSlice("volumes")
	if c.Bool("local") {
		workspace := filepath.Join(
			c.String("workspace-base"),
			c.String("workspace-path"),
		)

		dir, _ := filepath.Abs(filepath.Dir(file))
		volumes = append(volumes, dir+":"+workspace)
	}

	// compiles the yaml file
	compiled := bitbucket.NewCompiler(
		bitbucket.WithVolumes(volumes...),
		bitbucket.WithWorkspace(
			c.String("workspace-base"),
			c.String("workspace-path"),
		),
		bitbucket.WithPrefix(
			c.String("prefix"),
		),
		bitbucket.WithLocal(
			c.Bool("local"),
		),
		bitbucket.WithNetrc(
			c.String("netrc-username"),
			c.String("netrc-password"),
			c.String("netrc-machine"),
		),
		bitbucket.WithMetadata(
			metadataFromContext(c),
		),
	).Compile(conf)

	// marshal the compiled spec to formatted yaml
	out, err := json.MarshalIndent(compiled, "", "  ")
	if err != nil {
		return err
	}

	// create output file with option to dump to stdout
	var writer = os.Stdout
	output := c.String("out")
	if output != "-" {
		writer, err = os.Create(output)
		if err != nil {
			return err
		}
	}
	defer writer.Close()

	_, err = writer.Write(out)
	if err != nil {
		return err
	}

	if writer != os.Stdout {
		fmt.Fprintf(os.Stdout, "Successfully compiled %s to %s\n", file, output)
	}
	return nil
}

// return the metadata from the cli context.
func metadataFromContext(c *cli.Context) frontend.Metadata {
	return frontend.Metadata{
		Repo: frontend.Repo{
			Name:    c.String("repo-name"),
			Link:    c.String("repo-link"),
			Remote:  c.String("repo-remote-url"),
			Private: c.Bool("repo-private"),
		},
		Curr: frontend.Build{
			Number:   c.Int("build-number"),
			Created:  c.Int64("build-created"),
			Started:  c.Int64("build-started"),
			Finished: c.Int64("build-finished"),
			Status:   c.String("build-status"),
			Event:    c.String("build-event"),
			Link:     c.String("build-link"),
			Target:   c.String("build-target"),
			Commit: frontend.Commit{
				Sha:     c.String("commit-sha"),
				Ref:     c.String("commit-ref"),
				Refspec: c.String("commit-refspec"),
				Branch:  c.String("commit-branch"),
				Message: c.String("commit-message"),
				Author: frontend.Author{
					Name:   c.String("commit-author-name"),
					Email:  c.String("commit-author-email"),
					Avatar: c.String("commit-author-avatar"),
				},
			},
		},
		Prev: frontend.Build{
			Number:   c.Int("prev-build-number"),
			Created:  c.Int64("prev-build-created"),
			Started:  c.Int64("prev-build-started"),
			Finished: c.Int64("prev-build-finished"),
			Status:   c.String("prev-build-status"),
			Event:    c.String("prev-build-event"),
			Link:     c.String("prev-build-link"),
			Commit: frontend.Commit{
				Sha:     c.String("prev-commit-sha"),
				Ref:     c.String("prev-commit-ref"),
				Refspec: c.String("prev-commit-refspec"),
				Branch:  c.String("prev-commit-branch"),
				Message: c.String("prev-commit-message"),
				Author: frontend.Author{
					Name:   c.String("prev-commit-author-name"),
					Email:  c.String("prev-commit-author-email"),
					Avatar: c.String("prev-commit-author-avatar"),
				},
			},
		},
		Job: frontend.Job{
			Number: c.Int("job-number"),
			// Matrix: ,
		},
		Sys: frontend.System{
			Name: c.String("system-name"),
			Link: c.String("system-link"),
			Arch: c.String("system-arch"),
		},
	}
}
