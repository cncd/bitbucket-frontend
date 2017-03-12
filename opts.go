package bitbucket

import "github.com/cncd/pipeline/pipeline/frontend"

// Option configures a compiler option.
type Option func(*Compiler)

// WithVolumes configutes the compiler with default volumes that
// are mounted to each container in the pipeline.
func WithVolumes(volumes ...string) Option {
	return func(compiler *Compiler) {
		compiler.volumes = volumes
	}
}

// WithWorkspace configures the compiler with the workspace base
// and path. The workspace base is a volume created at runtime and
// mounted into all containers in the pipeline. The base and path
// are joined to provide the working directory for all build and
// plugin steps in the pipeline.
func WithWorkspace(base, path string) Option {
	return func(compiler *Compiler) {
		compiler.base = base
		compiler.path = path
	}
}

// WithMetadata configutes the compiler with the repostiory, build
// and system metadata. The metadata is used to remove steps from
// the compiled pipeline configuration that should be skipped. The
// metadata is also added to each container as environment variables.
func WithMetadata(metadata frontend.Metadata) Option {
	return func(compiler *Compiler) {
		compiler.meta = metadata
		for k, v := range metadata.Environ() {
			compiler.env[k] = v
		}
		// TODO deprecate these environment variables
		for k, v := range metadata.EnvironDrone() {
			compiler.env[k] = v
		}
	}
}

// WithNetrc configures the compiler with netrc authentication
// credentials added by default to every container in the pipeline.
func WithNetrc(username, password, machine string) Option {
	return WithEnviron(
		map[string]string{
			"CI_NETRC_USERNAME": username,
			"CI_NETRC_PASSWORD": password,
			"CI_NETRC_MACHINE":  machine,

			// TODO deprecate these environment variables
			"DRONE_NETRC_USERNAME": username,
			"DRONE_NETRC_PASSWORD": password,
			"DRONE_NETRC_MACHINE":  machine,
		},
	)
}

// WithPrefix configures the compiler with the prefix. The prefix is
// used to prefix container, volume and network names to avoid
// collision at runtime.
func WithPrefix(prefix string) Option {
	return func(compiler *Compiler) {
		compiler.prefix = prefix
	}
}

// WithEnviron configures the compiler with environment variables
// added by default to every container in the pipeline.
func WithEnviron(env map[string]string) Option {
	return func(compiler *Compiler) {
		for k, v := range env {
			compiler.env[k] = v
		}
	}
}

// WithLocal configures the compiler with the local flag. The local
// flag indicates the pipeline execution is running in a local development
// environment with a mounted local working directory.
func WithLocal(local bool) Option {
	return func(compiler *Compiler) {
		compiler.local = local
	}
}

// WithProxy configures the compiler with HTTP_PROXY, HTTPS_PROXY,
// and NO_PROXY environment variables added by default to every
// container in the pipeline.
func WithProxy(http, https, none string) Option {
	return WithEnviron(
		map[string]string{
			"no_proxy":    none,
			"NO_PROXY":    none,
			"http_proxy":  http,
			"HTTP_PROXY":  http,
			"HTTPS_PROXY": https,
			"https_proxy": https,
		},
	)
}
