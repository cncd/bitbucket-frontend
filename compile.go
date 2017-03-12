package bitbucket

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/cncd/pipeline/pipeline/backend"
	"github.com/cncd/pipeline/pipeline/frontend"
	"github.com/docker/docker/reference"
)

// Compiler compiles the yaml
type Compiler struct {
	local   bool
	prefix  string
	volumes []string
	env     map[string]string
	base    string
	path    string
	meta    frontend.Metadata
}

// NewCompiler ...
func NewCompiler(opts ...Option) *Compiler {
	compiler := new(Compiler)
	compiler.env = map[string]string{}
	compiler.base = "/workspace"
	compiler.path = "src"
	for _, opt := range opts {
		opt(compiler)
	}
	compiler.env["CI_WORKSPACE"] = path.Join(compiler.base, compiler.path)

	// legacy code for DRONE_ plugins
	for k, v := range compiler.env {
		if strings.HasPrefix(k, "CI_") {
			p := strings.Replace(k, "CI_", "DRONE_", 1)
			compiler.env[p] = v
		}
	}
	return compiler
}

// Compile compiles the Yaml to the common runtime.
func (c *Compiler) Compile(conf *Config) *backend.Config {
	spec := new(backend.Config)

	// choose which pipeline to execute
	// return the pipeline by name
	section := conf.Pipeline(c.meta.Curr.Commit.Ref, c.meta.Curr.Commit.Branch)

	// defines the default workspace
	workingdir := path.Join(c.base, c.path)

	// creates the default volume.
	volume := new(backend.Volume)
	volume.Driver = "local"
	volume.Name = fmt.Sprintf("%s_workspace", c.prefix)
	spec.Volumes = append(spec.Volumes, volume)

	// create the default volume reference.
	volumes := []string{
		volume.Name + ":" + c.base,
	}
	volumes = append(volumes, c.volumes...)

	// adds the default clone stage
	if c.local == false {
		envs := copyEnv(c.env)
		envs["PLUGIN_DEPTH"] = strconv.Itoa(conf.Clone.Depth)

		step := &backend.Step{
			Name:        fmt.Sprintf("%s_clone", c.prefix),
			Alias:       "clone",
			Image:       "plugins/git:latest",
			Environment: envs,
			OnSuccess:   true,
			OnFailure:   false,
			Volumes:     volumes,
			WorkingDir:  workingdir,
		}

		stage := new(backend.Stage)
		stage.Name = fmt.Sprintf("%s_clone", c.prefix)
		stage.Alias = "clone"
		stage.Steps = append(stage.Steps, step)

		spec.Stages = append(spec.Stages, stage)
	}

	// adds the pipeline steps
	for i, step := range section.Steps {
		image := step.Image
		if image == "" {
			image = conf.Image
		}
		image = expandImage(image)

		envs := copyEnv(c.env)
		envs["CI_SCRIPT"] = toScript(step.Script)
		envs["DRONE_SCRIPT"] = toScript(step.Script)
		envs["HOME"] = "/root"
		envs["SHELL"] = "/bin/sh"

		step := &backend.Step{
			Name:        fmt.Sprintf("%s_step_%d", c.prefix, i),
			Alias:       fmt.Sprintf("step_%d", i),
			Image:       image,
			Environment: envs,
			Entrypoint:  []string{"/bin/sh", "-c"},
			Command:     []string{"echo $CI_SCRIPT | base64 -d | /bin/sh -e"},
			Volumes:     volumes,
			WorkingDir:  workingdir,
			OnSuccess:   true,
			OnFailure:   false,
		}

		stage := new(backend.Stage)
		stage.Name = fmt.Sprintf("%s_stage_%d", c.prefix, i)
		stage.Alias = fmt.Sprintf("stage_%d", i)
		stage.Steps = append(stage.Steps, step)

		spec.Stages = append(spec.Stages, stage)
	}

	return spec
}

func copyEnv(from map[string]string) map[string]string {
	to := map[string]string{}
	for k, v := range from {
		to[k] = v
	}
	return to
}

func expandImage(name string) string {
	ref, err := reference.ParseNamed(name)
	if err != nil {
		return name
	}
	return reference.WithDefaultTag(ref).String()
}

func toScript(commands []string) string {
	var buf bytes.Buffer
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}

	script := fmt.Sprintf(
		setupScript,
		buf.String(),
	)

	return base64.StdEncoding.EncodeToString([]byte(script))
}

// setupScript is a helper script this is added to the build to ensure
// a minimum set of environment variables are set correctly.
const setupScript = `
if [ -n "$CI_NETRC_MACHINE" ]; then
cat <<EOF > $HOME/.netrc
machine $CI_NETRC_MACHINE
login $CI_NETRC_USERNAME
password $CI_NETRC_PASSWORD
EOF
chmod 0600 $HOME/.netrc
fi
unset CI_NETRC_USERNAME
unset CI_NETRC_PASSWORD
unset CI_SCRIPT
%s
`

// traceScript is a helper script that is added to the build script
// to trace a command.
const traceScript = `
echo + %s
%s
`
