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
			Name:    "build-args",
			Usage:   "build args",
			EnvVars: []string{"PLUGIN_BUILD_ARGS"},
		},
		&cli.BoolFlag{
			Name:    "cache",
			Usage:   "enable kaniko caching",
			EnvVars: []string{"PLUGIN_CACHE"},
		},
		&cli.StringFlag{
			Name:    "cache-dir",
			Usage:   "cache dir",
			EnvVars: []string{"PLUGIN_CACHE_DIR"},
		},
		&cli.StringFlag{
			Name:    "cache-repo",
			Usage:   "cache repository",
			EnvVars: []string{"PLUGIN_CACHE_REPO"},
		},
		&cli.DurationFlag{
			Name:    "cache-ttl",
			Usage:   "cache ttl",
			EnvVars: []string{"PLUGIN_CACHE_TTL"},
		},
		&cli.BoolFlag{
			Name:    "cleanup",
			Usage:   "cleanup",
			EnvVars: []string{"PLUGIN_CLEANUP", "PLUGIN_PURGE"},
		},
		&cli.StringFlag{
			Name:    "context",
			Usage:   "build context",
			EnvVars: []string{"PLUGIN_CONTEXT"},
		},
		&cli.StringSliceFlag{
			Name:    "destinations",
			Usage:   "destinations",
			EnvVars: []string{"PLUGIN_DESTINATIONS"},
		},
		&cli.StringFlag{
			Name:    "digest-file",
			Usage:   "digest file",
			EnvVars: []string{"PLUGIN_DIGEST_FILE"},
		},
		&cli.StringFlag{
			Name:    "dockerfile",
			Usage:   "build dockerfile",
			EnvVars: []string{"PLUGIN_DOCKERFILE"},
		},
		&cli.BoolFlag{
			Name:    "force",
			Usage:   "force",
			EnvVars: []string{"PLUGIN_FORCE"},
		},
		&cli.StringFlag{
			Name:    "image-name-with-digest-file",
			Usage:   "image name with digest file",
			EnvVars: []string{"PLUGIN_IMAGE_NAME_WITH_DIGEST_FILE"},
		},
		&cli.BoolFlag{
			Name:    "insecure-push",
			Usage:   "insecure push",
			EnvVars: []string{"PLUGIN_INSECURE_PUSH"},
		},
		&cli.BoolFlag{
			Name:    "insecure-pull",
			Usage:   "insecure pull",
			EnvVars: []string{"PLUGIN_INSECURE_PULL"},
		},
		&cli.StringSliceFlag{
			Name:    "insecure-registries",
			Usage:   "insecure registry",
			EnvVars: []string{"PLUGIN_INSECURE_REGISTRY"},
		},
		&cli.StringSliceFlag{
			Name:    "labels",
			Usage:   "labels",
			EnvVars: []string{"PLUGIN_LABELS", "PLUGIN_CUSTOM_LABELS"},
		},
		&cli.StringFlag{
			Name:    "log-format",
			Usage:   "log format",
			EnvVars: []string{"PLUGIN_LOG_FORMAT"},
		},
		&cli.BoolFlag{
			Name:    "no-push",
			Usage:   "no push",
			EnvVars: []string{"PLUGIN_NO_PUSH", "PLUGIN_DRY_RUN"},
		},
		&cli.StringFlag{
			Name:    "oci-layout-path",
			Usage:   "oci layout path",
			EnvVars: []string{"PLUGIN_OCI_LAYOUT_PATH"},
		},
		&cli.StringFlag{
			Name:    "registry-mirror",
			Usage:   "registry mirror",
			EnvVars: []string{"PLUGIN_REGISTRY_MIRROR", "PLUGIN_MIRROR"},
		},
		&cli.BoolFlag{
			Name:    "reproducible",
			Usage:   "reproducible",
			EnvVars: []string{"PLUGIN_REPRODUCIBLE"},
		},
		&cli.BoolFlag{
			Name:    "single-snapshot",
			Usage:   "single snapshot",
			EnvVars: []string{"PLUGIN_SINGLE_SNAPSHOT"},
		},
		&cli.BoolFlag{
			Name:    "skip-tls-verify",
			Usage:   "skip tls verify",
			EnvVars: []string{"PLUGIN_SKIP_TLS_VERIFY"},
		},
		&cli.BoolFlag{
			Name:    "skip-tls-verify-pull",
			Usage:   "skip tls verify-pull",
			EnvVars: []string{"PLUGIN_SKIP_TLS_VERIFY_PULL"},
		},
		&cli.StringSliceFlag{
			Name:    "skip-tls-verify-registries",
			Usage:   "skip tls verify registry",
			EnvVars: []string{"PLUGIN_SKIP_TLS_VERIFY_REGISTRY"},
		},
		&cli.StringFlag{
			Name:    "snapshot-mode",
			Usage:   "snapshot mode",
			EnvVars: []string{"PLUGIN_SNAPSHOT_MODE"},
		},
		&cli.StringFlag{
			Name:    "tar-path",
			Usage:   "tar path",
			EnvVars: []string{"PLUGIN_TAR_PATH"},
		},
		&cli.StringFlag{
			Name:    "target",
			Usage:   "target",
			EnvVars: []string{"PLUGIN_TARGET"},
		},
		&cli.StringFlag{
			Name:    "verbosity",
			Usage:   "verbosity",
			EnvVars: []string{"PLUGIN_VERBOSITY"},
		},
		&cli.BoolFlag{
			Name:    "whitelist-var-run",
			Usage:   "whitelist /var/run/*",
			EnvVars: []string{"PLUGIN_WHITELIST_VAR_RUN"},
		},
		// docker auth args
		&cli.StringFlag{
			Name:    "docker.registry",
			Usage:   "docker registry",
			EnvVars: []string{"PLUGIN_REGISTRY"},
		},
		&cli.StringFlag{
			Name:    "docker.username",
			Usage:   "docker username",
			EnvVars: []string{"PLUGIN_USERNAME", "DOCKER_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "docker.password",
			Usage:   "docker password",
			EnvVars: []string{"PLUGIN_PASSWORD", "DOCKER_PASSWORD"},
		},
		// other flags
		&cli.StringFlag{
			Name:    "build.repo",
			Usage:   "repository name",
			EnvVars: []string{"PLUGIN_REPO"},
		},
		&cli.BoolFlag{
			Name:    "build.debug",
			Usage:   "enable debug mode",
			EnvVars: []string{"PLUGIN_DEBUG", "DOCKER_LAUNCH_DEBUG"},
		},
		&cli.StringSliceFlag{
			Name:     "build.tags",
			Usage:    "build tags",
			Value:    cli.NewStringSlice("latest"),
			EnvVars:  []string{"PLUGIN_TAG", "PLUGIN_TAGS"},
			FilePath: ".tags",
		},
		&cli.StringSliceFlag{
			Name:    "build.images",
			Usage:   "cache from repo",
			EnvVars: []string{"PLUGIN_WARMER_IMAGES", "PLUGIN_CACHE_FROM"},
		},
		&cli.BoolFlag{
			Name:    "build.force-cache",
			Usage:   "Force cache overwritting",
			EnvVars: []string{"PLUGIN_FORCE_CACHE"},
		},
		&cli.StringSliceFlag{
			Name:    "build.args-from-env",
			Usage:   "build args",
			EnvVars: []string{"PLUGIN_BUILD_ARGS_FROM_ENV"},
		},
		&cli.BoolFlag{
			Name:    "build.tags-auto",
			Usage:   "default build tags",
			EnvVars: []string{"PLUGIN_DEFAULT_TAGS,PLUGIN_AUTO_TAG"},
		},
		&cli.StringFlag{
			Name:    "build.tags-suffix",
			Usage:   "default build tags with suffix",
			EnvVars: []string{"PLUGIN_DEFAULT_SUFFIX", "PLUGIN_AUTO_TAG_SUFFIX"},
		},
		&cli.BoolFlag{
			Name:    "build.insecure",
			Usage:   "insecure",
			EnvVars: []string{"PLUGIN_INSECURE"},
		},
	}
}

// settingsFromContext creates a plugin.Settings from the cli.Context.
func settingsFromContext(ctx *cli.Context) kaniko.Settings {
	return kaniko.Settings{
		BuildArgs:               ctx.StringSlice("build-args"),
		Cache:                   ctx.Bool("cache"),
		CacheDir:                ctx.String("cache-dir"),
		CacheRepo:               ctx.String("cache-repo"),
		CacheTTL:                ctx.Duration("cache-ttl"),
		Cleanup:                 ctx.Bool("cleanup"),
		Context:                 ctx.String("context"),
		Destinations:            ctx.StringSlice("destinations"),
		DigestFile:              ctx.String("digest-file"),
		Dockerfile:              ctx.String("dockerfile"),
		Force:                   ctx.Bool("force"),
		ImageNameWithDigestFile: ctx.String("image-name-with-digest-file"),
		InsecurePush:            ctx.Bool("insecure-push"),
		InsecurePull:            ctx.Bool("insecure-pull"),
		InsecureRegistries:      ctx.StringSlice("insecure-registries"),
		Labels:                  ctx.StringSlice("labels"),
		LogFormat:               ctx.String("log-format"),
		NoPush:                  ctx.Bool("no-push"),
		OCILayoutPath:           ctx.String("oci-layout-path"),
		RegistryCertificates:    ctx.StringSlice("registry-certificates"),
		RegistryMirror:          ctx.String("registry-mirror"),
		Reproducible:            ctx.Bool("reproducible"),
		SingleSnapshot:          ctx.Bool("single-snapshot"),
		SkipTLSVerify:           ctx.Bool("skip-tls-verify"),
		SkipTLSVerifyPull:       ctx.Bool("skip-tls-verify-pull"),
		SkipTLSVerifyRegistries: ctx.StringSlice("skip-tls-verify-registries"),
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
		BuildArgsFromEnv: ctx.StringSlice("build.args-from-env"),
		ForceCache:       ctx.Bool("build.force-cache"),
		Tags:             ctx.StringSlice("build.tags"),
		TagsAuto:         ctx.Bool("build.tags-auto"),
		TagsSuffix:       ctx.String("build.tags-suffix"),
		WarmerImages:     ctx.StringSlice("build.images"),
		Repo:             ctx.String("build.repo"),
		Debug:            ctx.Bool("build.debug"),
		Insecure:         ctx.Bool("build.insecure"),
	}
}
