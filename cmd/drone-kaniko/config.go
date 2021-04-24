package main

import (
	"github.com/urfave/cli/v2"
	"megpoid.xyz/go/drone-kaniko/pkg/kaniko"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags() []cli.Flag {
	// Replace below with all the flags required for the plugin's specific
	// settings.
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "build-arg",
			Aliases: []string{"args"},
			Usage:   `This flag allows you to pass in ARG values at build time`,
			EnvVars: []string{"PLUGIN_BUILD_ARGS"},
		},
		&cli.BoolFlag{
			Name:    "cache",
			Usage:   `Use cache when building image`,
			EnvVars: []string{"PLUGIN_CACHE"},
		},
		&cli.StringFlag{
			Name:    "cache-dir",
			Usage:   `Specify a local directory to use as a cache")`,
			EnvVars: []string{"PLUGIN_CACHE_DIR"},
		},
		&cli.StringFlag{
			Name:    "cache-repo",
			Usage:   `Specify a repository to use as a cache, otherwise one will be inferred from the destination provided`,
			EnvVars: []string{"PLUGIN_CACHE_REPO"},
		},
		&cli.DurationFlag{
			Name:    "cache-ttl",
			Usage:   `Cache timeout in hours`,
			EnvVars: []string{"PLUGIN_CACHE_TTL"},
		},
		&cli.BoolFlag{
			Name:    "cleanup",
			Aliases: []string{"purge"},
			Usage:   `Clean the filesystem at the end`,
			EnvVars: []string{"PLUGIN_CLEANUP", "PLUGIN_PURGE"},
		},
		&cli.StringFlag{
			Name:    "context",
			Usage:   `Path to the dockerfile build context`,
			EnvVars: []string{"PLUGIN_CONTEXT"},
		},
		&cli.StringSliceFlag{
			Name:    "destination",
			Usage:   `Registry the final image should be pushed to`,
			EnvVars: []string{"PLUGIN_DESTINATIONS"},
		},
		&cli.StringFlag{
			Name:    "digest-file",
			Usage:   `Specify a file to save the digest of the built image to`,
			EnvVars: []string{"PLUGIN_DIGEST_FILE"},
		},
		&cli.StringFlag{
			Name:    "dockerfile",
			Usage:   `Path to the dockerfile to be built`,
			EnvVars: []string{"PLUGIN_DOCKERFILE"},
		},
		&cli.BoolFlag{
			Name:    "force",
			Usage:   `Force building outside of a container`,
			EnvVars: []string{"PLUGIN_FORCE"},
		},
		&cli.StringFlag{
			Name:    "image-name-with-digest-file",
			Usage:   `Specify a file to save the image name w/ digest of the built image to`,
			EnvVars: []string{"PLUGIN_IMAGE_NAME_WITH_DIGEST_FILE"},
		},
		&cli.BoolFlag{
			Name:    "insecure-push",
			Usage:   `Push to insecure registry using plain HTTP`,
			EnvVars: []string{"PLUGIN_INSECURE_PUSH"},
		},
		&cli.BoolFlag{
			Name:    "insecure-pull",
			Usage:   `Pull from insecure registry using plain HTTP`,
			EnvVars: []string{"PLUGIN_INSECURE_PULL"},
		},
		&cli.StringSliceFlag{
			Name:    "insecure-registry",
			Usage:   `Insecure registry using plain HTTP to push and pull`,
			EnvVars: []string{"PLUGIN_INSECURE_REGISTRY"},
		},
		&cli.StringSliceFlag{
			Name:    "label",
			Usage:   `Set metadata for an image`,
			EnvVars: []string{"PLUGIN_LABELS", "PLUGIN_CUSTOM_LABELS"},
		},
		&cli.StringFlag{
			Name:    "log-format",
			Usage:   `Log format (text, color, json)`,
			EnvVars: []string{"PLUGIN_LOG_FORMAT"},
		},
		&cli.BoolFlag{
			Name:    "no-push",
			Aliases: []string{"dry-run"},
			Usage:   `Do not push the image to the registry`,
			EnvVars: []string{"PLUGIN_NO_PUSH", "PLUGIN_DRY_RUN"},
		},
		&cli.StringFlag{
			Name:    "oci-layout-path",
			Usage:   `Path to save the OCI image layout of the built image`,
			EnvVars: []string{"PLUGIN_OCI_LAYOUT_PATH"},
		},
		&cli.StringSliceFlag{
			Name:    "registry-certificate",
			Usage:   `Use the provided certificate for TLS communication with the given registry. Expected format is 'my.registry.url=/path/to/the/server/certificate'`,
			EnvVars: []string{"PLUGIN_REGISTRY_CERTIFICATES"},
		},
		&cli.StringFlag{
			Name:    "registry-mirror",
			Usage:   `Registry mirror to use has pull-through cache instead of docker.io`,
			EnvVars: []string{"PLUGIN_REGISTRY_MIRROR"},
		},
		&cli.BoolFlag{
			Name:    "reproducible",
			Usage:   `Strip timestamps out of the image to make it reproducible`,
			EnvVars: []string{"PLUGIN_REPRODUCIBLE"},
		},
		&cli.BoolFlag{
			Name:    "single-snapshot",
			Usage:   `Take a single snapshot at the end of the build`,
			EnvVars: []string{"PLUGIN_SINGLE_SNAPSHOT"},
		},
		&cli.BoolFlag{
			Name:    "skip-tls-verify",
			Usage:   `Push to insecure registry ignoring TLS verify`,
			EnvVars: []string{"PLUGIN_SKIP_TLS_VERIFY"},
		},
		&cli.BoolFlag{
			Name:    "skip-tls-verify-pull",
			Usage:   `Pull from insecure registry ignoring TLS verify`,
			EnvVars: []string{"PLUGIN_SKIP_TLS_VERIFY_PULL"},
		},
		&cli.StringSliceFlag{
			Name:    "skip-tls-verify-registry",
			Usage:   `Insecure registry ignoring TLS verify to push and pull`,
			EnvVars: []string{"PLUGIN_SKIP_TLS_VERIFY_REGISTRY"},
		},
		&cli.StringFlag{
			Name:    "snapshot-mode",
			Usage:   `Change the file attributes inspected during snapshotting`,
			EnvVars: []string{"PLUGIN_SNAPSHOT_MODE"},
		},
		&cli.StringFlag{
			Name:    "tar-path",
			Usage:   `Path to save the image in as a tarball instead of pushing`,
			EnvVars: []string{"PLUGIN_TAR_PATH"},
		},
		&cli.StringFlag{
			Name:    "target",
			Usage:   `Set the target build stage to build`,
			EnvVars: []string{"PLUGIN_TARGET"},
		},
		&cli.StringFlag{
			Name:    "verbosity",
			Usage:   `Log level (debug, info, warn, error, fatal, panic`,
			EnvVars: []string{"PLUGIN_VERBOSITY"},
		},
		&cli.BoolFlag{
			Name:    "whitelist-var-run",
			Usage:   `Ignore /var/run directory when taking image snapshot. Set it to false to preserve /var/run/ in destination image`,
			EnvVars: []string{"PLUGIN_WHITELIST_VAR_RUN"},
		},
		// docker auth args
		&cli.StringFlag{
			Name:    "docker.registry",
			Usage:   `Docker registry`,
			Value:   "https://index.docker.io/v1/",
			EnvVars: []string{"PLUGIN_REGISTRY"},
		},
		&cli.StringFlag{
			Name:    "docker.username",
			Usage:   `Docker username`,
			EnvVars: []string{"PLUGIN_USERNAME", "DOCKER_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "docker.password",
			Usage:   `Docker password`,
			EnvVars: []string{"PLUGIN_PASSWORD", "DOCKER_PASSWORD"},
		},
		// other flags
		&cli.StringFlag{
			Name:    "repo",
			Usage:   `Docker repository name. Compatible with drone-docker plugin to provide 'destination'`,
			EnvVars: []string{"PLUGIN_REPO"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   `Enable debug logging. Compatible with drone-docker plugin to provide 'verbosity'`,
			EnvVars: []string{"PLUGIN_DEBUG", "DOCKER_LAUNCH_DEBUG"},
		},
		&cli.StringSliceFlag{
			Name:     "tags",
			Usage:    `Build tags. Compatible with drone-docker plugin to provide 'destination'`,
			Value:    cli.NewStringSlice("latest"),
			EnvVars:  []string{"PLUGIN_TAG", "PLUGIN_TAGS"},
			FilePath: ".tags",
		},
		&cli.StringSliceFlag{
			Name:    "image",
			Aliases: []string{"cache-from"},
			Usage:   `Image to cache`,
			EnvVars: []string{"PLUGIN_WARMER_IMAGES", "PLUGIN_CACHE_FROM"},
		},
		&cli.BoolFlag{
			Name:    "force-cache",
			Usage:   `Force cache overwritting`,
			EnvVars: []string{"PLUGIN_FORCE_CACHE"},
		},
		&cli.StringSliceFlag{
			Name:    "args-from-env",
			Usage:   "This flag allows you to pass in ARG values at build time. Read from environment",
			EnvVars: []string{"PLUGIN_BUILD_ARGS_FROM_ENV"},
		},
		&cli.BoolFlag{
			Name:    "tags-auto",
			Usage:   `Default build tags`,
			EnvVars: []string{"PLUGIN_DEFAULT_TAGS,PLUGIN_AUTO_TAG"},
		},
		&cli.StringFlag{
			Name:    "tags-suffix",
			Usage:   `Default build tags with suffix`,
			EnvVars: []string{"PLUGIN_DEFAULT_SUFFIX", "PLUGIN_AUTO_TAG_SUFFIX"},
		},
		&cli.BoolFlag{
			Name:    "insecure",
			Usage:   `Allow insecure registries`,
			EnvVars: []string{"PLUGIN_INSECURE"},
		},
		&cli.StringSliceFlag{
			Name:    "label-schema",
			Usage:   `Label-schema labels`,
			EnvVars: []string{"PLUGIN_LABEL_SCHEMA"},
		},
		&cli.StringFlag{
			Name:    "mirror",
			Usage:   `Registry mirror (e.g. https://mirror.example.com). Compatible with drone-docker plugin to provide 'registry-mirror'`,
			EnvVars: []string{"PLUGIN_MIRROR"},
		},
	}
}

// settingsFromContext creates a plugin.Settings from the cli.Context.
func settingsFromContext(ctx *cli.Context) kaniko.Settings {
	return kaniko.Settings{
		BuildArgs:               ctx.StringSlice("build-arg"),
		Cache:                   ctx.Bool("cache"),
		CacheDir:                ctx.String("cache-dir"),
		CacheRepo:               ctx.String("cache-repo"),
		CacheTTL:                ctx.Duration("cache-ttl"),
		Cleanup:                 ctx.Bool("cleanup"),
		Context:                 ctx.String("context"),
		Destinations:            ctx.StringSlice("destination"),
		DigestFile:              ctx.String("digest-file"),
		Dockerfile:              ctx.String("dockerfile"),
		Force:                   ctx.Bool("force"),
		ImageNameWithDigestFile: ctx.String("image-name-with-digest-file"),
		InsecurePush:            ctx.Bool("insecure-push"),
		InsecurePull:            ctx.Bool("insecure-pull"),
		InsecureRegistries:      ctx.StringSlice("insecure-registry"),
		Labels:                  ctx.StringSlice("label"),
		LogFormat:               ctx.String("log-format"),
		NoPush:                  ctx.Bool("no-push"),
		OCILayoutPath:           ctx.String("oci-layout-path"),
		RegistryCertificates:    ctx.StringSlice("registry-certificate"),
		RegistryMirror:          ctx.String("registry-mirror"),
		Reproducible:            ctx.Bool("reproducible"),
		SingleSnapshot:          ctx.Bool("single-snapshot"),
		SkipTLSVerify:           ctx.Bool("skip-tls-verify"),
		SkipTLSVerifyPull:       ctx.Bool("skip-tls-verify-pull"),
		SkipTLSVerifyRegistries: ctx.StringSlice("skip-tls-verify-registry"),
		SnapshotMode:            ctx.String("snapshot-mode"),
		TarPath:                 ctx.String("tar-path"),
		Target:                  ctx.String("target"),
		Verbosity:               ctx.String("verbosity"),
		WhitelistVarRun:         ctx.Bool("whitelist-var-run"),
		// auth args
		Registry: ctx.String("docker.registry"),
		Username: ctx.String("docker.username"),
		Password: ctx.String("docker.password"),
		// other args
		BuildArgsFromEnv: ctx.StringSlice("args-from-env"),
		ForceCache:       ctx.Bool("force-cache"),
		Tags:             ctx.StringSlice("tags"),
		TagsAuto:         ctx.Bool("tags-auto"),
		TagsSuffix:       ctx.String("tags-suffix"),
		Images:           ctx.StringSlice("image"),
		Repo:             ctx.String("repo"),
		Debug:            ctx.Bool("debug"),
		Insecure:         ctx.Bool("insecure"),
		LabelSchema:      ctx.StringSlice("label-schema"),
		Mirror:           ctx.String("mirror"),
		PushTarget:       ctx.Bool("push-target"),
	}
}
