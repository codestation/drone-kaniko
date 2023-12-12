package main

import (
	"github.com/urfave/cli/v2"
	"go.megpoid.dev/drone-kaniko/pkg/kaniko"
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
		&cli.BoolFlag{
			Name:    "cache-copy-layers",
			Usage:   `Caches copy layers`,
			EnvVars: []string{"PLUGIN_CACHE_COPY_LAYERS"},
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
		&cli.BoolFlag{
			Name:    "compressed-caching",
			Usage:   `Compress the cached layers. Decreases build time, but increases memory usage`,
			EnvVars: []string{"PLUGIN_COMPRESSED_CACHING"},
			Value:   true,
		},
		&cli.StringFlag{
			Name:    "compression",
			Usage:   `Compression algorithm (gzip, zstd)`,
			EnvVars: []string{"PLUGIN_COMPRESSION"},
		},
		&cli.IntFlag{
			Name:    "compression-level",
			Usage:   `Compression level`,
			EnvVars: []string{"PLUGIN_COMPRESSION_LEVEL"},
			Value:   -1,
		},
		&cli.StringFlag{
			Name:    "context",
			Usage:   `Path to the dockerfile build context`,
			EnvVars: []string{"PLUGIN_CONTEXT"},
		},
		&cli.StringFlag{
			Name:    "context-sub-path",
			Usage:   `Sub path within the given context`,
			EnvVars: []string{"PLUGIN_CONTEXT_SUB_PATH"},
		},
		&cli.StringFlag{
			Name:    "custom-platform",
			Usage:   `Specify the build platform if different from the current host`,
			EnvVars: []string{"PLUGIN_CUSTOM_PLATFORM"},
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
		&cli.BoolFlag{
			Name:    "force-build-metadata",
			Usage:   `Force add metadata layers to build image`,
			EnvVars: []string{"PLUGIN_FORCE_BUILD_METADATA"},
		},
		&cli.StringFlag{
			Name:    "git",
			Usage:   `Branch to clone if build context is a git repository`,
			EnvVars: []string{"PLUGIN_GIT"},
		},
		&cli.StringSliceFlag{
			Name:    "ignore-path",
			Usage:   `Ignore these paths when taking a snapshot`,
			EnvVars: []string{"PLUGIN_IGNORE_PATH"},
		},
		&cli.BoolFlag{
			Name:    "ignore-var-run",
			Usage:   `Ignore /var/run directory when taking image snapshot`,
			EnvVars: []string{"PLUGIN_IGNORE_VAR_RUN"},
			Value:   true,
		},
		&cli.IntFlag{
			Name:    "image-fs-extract-retry",
			Usage:   `Number of retries for image FS extraction`,
			EnvVars: []string{"PLUGIN_IMAGE_FS_EXTRACT_RETRY"},
		},
		&cli.StringFlag{
			Name:    "image-name-tag-with-digest-file",
			Usage:   `Specify a file to save the image name w/ image tag w/ digest of the built image to`,
			EnvVars: []string{"PLUGIN_IMAGE_NAME_TAG_WITH_DIGEST_FILE"},
		},
		&cli.StringFlag{
			Name:    "image-name-with-digest-file",
			Usage:   `Specify a file to save the image name w/ digest of the built image to`,
			EnvVars: []string{"PLUGIN_IMAGE_NAME_WITH_DIGEST_FILE"},
		},
		&cli.BoolFlag{
			Name:    "insecure",
			Usage:   `Push to insecure registry using plain HTTP`,
			EnvVars: []string{"PLUGIN_INSECURE"},
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
		&cli.StringFlag{
			Name:    "kaniko-dir",
			Usage:   `Path to the kaniko directory, this takes precedence over the KANIKO_DIR environment variable`,
			EnvVars: []string{"PLUGIN_KANIKO_DIR"},
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
			Name:    "log-timestamp",
			Usage:   `Timestamp in log output`,
			EnvVars: []string{"PLUGIN_LOG_TIMESTAMP"},
		},
		&cli.BoolFlag{
			Name:    "no-push",
			Aliases: []string{"dry-run"},
			Usage:   `Do not push the image to the registry`,
			EnvVars: []string{"PLUGIN_NO_PUSH", "PLUGIN_DRY_RUN"},
		},
		&cli.BoolFlag{
			Name:    "no-push-cache",
			Usage:   `Do not push the cache layers to the registry`,
			EnvVars: []string{"PLUGIN_NO_PUSH_CACHE"},
		},
		&cli.StringFlag{
			Name:    "oci-layout-path",
			Usage:   `Path to save the OCI image layout of the built image`,
			EnvVars: []string{"PLUGIN_OCI_LAYOUT_PATH"},
		},
		&cli.IntFlag{
			Name:    "push-retry",
			Usage:   `Number of retries for the push operation`,
			EnvVars: []string{"PLUGIN_PUSH_RETRY"},
		},
		&cli.StringSliceFlag{
			Name:    "registry-certificate",
			Usage:   `Use the provided certificate for TLS communication with the given registry. Expected format is 'my.registry.url=/path/to/the/server/certificate'`,
			EnvVars: []string{"PLUGIN_REGISTRY_CERTIFICATES"},
		},
		&cli.StringSliceFlag{
			Name:    "registry-client-cert",
			Usage:   `Use the provided client certificate for mutual TLS (mTLS) communication with the given registry. Expected format is 'my.registry.url=/path/to/client/cert,/path/to/client/key''`,
			EnvVars: []string{"PLUGIN_REGISTRY_CLIENT_CERTS"},
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
			Name:    "skip-default-registry-fallback",
			Usage:   `If an image is not found on any mirrors (defined with registry-mirror) do not fallback to the default registry`,
			EnvVars: []string{"PLUGIN_SKIP_DEFAULT_REGISTRY_FALLBACK"},
		},
		&cli.BoolFlag{
			Name:    "skip-push-permission-check",
			Usage:   `Skip check of the push permission`,
			EnvVars: []string{"PLUGIN_SKIP_PUSH_PERMISSION_CHECK"},
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
		&cli.BoolFlag{
			Name:    "skip-unused-stages",
			Usage:   `Build only used stages if defined to true`,
			EnvVars: []string{"PLUGIN_SKIP_UNUSED_STAGES"},
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
		&cli.BoolFlag{
			Name:    "use-new-run",
			Usage:   `Use the experimental run implementation for detecting changes without requiring file system snapshots`,
			EnvVars: []string{"PLUGIN_USE_NEW_RUN"},
		},
		&cli.StringFlag{
			Name:    "verbosity",
			Usage:   `Log level (debug, info, warn, error, fatal, panic`,
			EnvVars: []string{"PLUGIN_VERBOSITY"},
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
		&cli.StringFlag{
			Name:    "docker.config",
			Usage:   `Docker config json content`,
			EnvVars: []string{"PLUGIN_CONFIG", "DOCKER_PLUGIN_CONFIG"},
		},
		// main flags
		&cli.StringSliceFlag{
			Name:    "args-from-env",
			Usage:   "This flag allows you to pass in ARG values at build time. Read from environment",
			EnvVars: []string{"PLUGIN_BUILD_ARGS_FROM_ENV"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   `Enable debug logging. Compatible with drone-docker plugin to provide 'verbosity'`,
			EnvVars: []string{"PLUGIN_DEBUG", "DOCKER_LAUNCH_DEBUG"},
		},
		&cli.StringFlag{
			Name:    "repo",
			Usage:   `Docker repository name. Compatible with drone-docker plugin to provide 'destination'`,
			EnvVars: []string{"PLUGIN_REPO"},
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
			Usage:   `Force cache overwritting (warmer)`,
			EnvVars: []string{"PLUGIN_FORCE_CACHE"},
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
		&cli.BoolFlag{
			Name:    "auto-label",
			Usage:   `Auto-label true|false`,
			Value:   true,
			EnvVars: []string{"PLUGIN_AUTO_LABEL"},
		},
		&cli.StringSliceFlag{
			Name:    "executor-extra-args",
			Usage:   "List of extra args to pass to the Kaniko executor process",
			EnvVars: []string{"PLUGIN_EXECUTOR_EXTRA_ARGS"},
		},
		&cli.StringSliceFlag{
			Name:    "warmer-extra-args",
			Usage:   "List of extra args to pass to the Kaniko warmer process",
			EnvVars: []string{"PLUGIN_WARMER_EXTRA_ARGS"},
		},
	}
}

// settingsFromContext creates a plugin.Settings from the cli.Context.
func settingsFromContext(ctx *cli.Context) kaniko.Settings {
	return kaniko.Settings{
		BuildArgs:                   ctx.StringSlice("build-arg"),
		Cache:                       ctx.Bool("cache"),
		CacheCopyLayers:             ctx.Bool("cache-copy-layers"),
		CacheDir:                    ctx.String("cache-dir"),
		CacheRepo:                   ctx.String("cache-repo"),
		CacheTTL:                    ctx.Duration("cache-ttl"),
		Cleanup:                     ctx.Bool("cleanup"),
		CompressedCaching:           ctx.Bool("compressed-caching"),
		Compression:                 ctx.String("compression"),
		CompressionLevel:            ctx.Int("compression-level"),
		Context:                     ctx.String("context"),
		ContextSubPath:              ctx.String("context-sub-path"),
		CustomPlatform:              ctx.String("custom-platform"),
		Destinations:                ctx.StringSlice("destination"),
		DigestFile:                  ctx.String("digest-file"),
		Dockerfile:                  ctx.String("dockerfile"),
		Force:                       ctx.Bool("force"),
		ForceBuildMetadata:          ctx.Bool("force-build-metadata"),
		Git:                         ctx.String("git"),
		IgnorePath:                  ctx.StringSlice("ignore-path"),
		IgnoreVarRun:                ctx.Bool("ignore-var-run"),
		ImageFsExtractRetry:         ctx.Int("image-fs-extract-retry"),
		ImageNameTagWithDigestFile:  ctx.String("image-name-tag-with-digest-file"),
		ImageNameWithDigestFile:     ctx.String("image-name-with-digest-file"),
		Insecure:                    ctx.Bool("insecure"),
		InsecurePull:                ctx.Bool("insecure-pull"),
		InsecureRegistries:          ctx.StringSlice("insecure-registry"),
		KanikoDir:                   ctx.String("kaniko-dir"),
		Labels:                      ctx.StringSlice("label"),
		LogFormat:                   ctx.String("log-format"),
		LogTimestamp:                ctx.Bool("log-timestamp"),
		NoPush:                      ctx.Bool("no-push"),
		NoPushCache:                 ctx.Bool("no-push-cache"),
		OCILayoutPath:               ctx.String("oci-layout-path"),
		PushRetry:                   ctx.Int("push-retry"),
		RegistryCertificates:        ctx.StringSlice("registry-certificate"),
		RegistryClientCerts:         ctx.StringSlice("registry-client-cert"),
		RegistryMirror:              ctx.String("registry-mirror"),
		Reproducible:                ctx.Bool("reproducible"),
		SingleSnapshot:              ctx.Bool("single-snapshot"),
		SkipDefaultRegistryFallback: ctx.Bool("skip-default-registry-fallback"),
		SkipPushPermissionCheck:     ctx.Bool("skip-push-permission-check"),
		SkipTLSVerify:               ctx.Bool("skip-tls-verify"),
		SkipTLSVerifyPull:           ctx.Bool("skip-tls-verify-pull"),
		SkipTLSVerifyRegistries:     ctx.StringSlice("skip-tls-verify-registry"),
		SkipUnusedStages:            ctx.Bool("skip-unused-stages"),
		SnapshotMode:                ctx.String("snapshot-mode"),
		TarPath:                     ctx.String("tar-path"),
		Target:                      ctx.String("target"),
		UseNewRun:                   ctx.Bool("use-new-run"),
		Verbosity:                   ctx.String("verbosity"),
		// auth args
		Auth: kaniko.Auth{
			Registry: ctx.String("docker.registry"),
			Username: ctx.String("docker.username"),
			Password: ctx.String("docker.password"),
			Config:   ctx.String("docker.config"),
		},
		// other args
		Main: kaniko.Main{
			BuildArgsFromEnv: ctx.StringSlice("args-from-env"),
			Debug:            ctx.Bool("debug"),
			DryRun:           ctx.Bool("dry-run"),
			ForceCache:       ctx.Bool("force-cache"),
			Tags:             ctx.StringSlice("tags"),
			TagsAuto:         ctx.Bool("tags-auto"),
			TagsSuffix:       ctx.String("tags-suffix"),
			Images:           ctx.StringSlice("image"),
			Repo:             ctx.String("repo"),
			LabelSchema:      ctx.StringSlice("label-schema"),
			Mirror:           ctx.String("mirror"),
			PushTarget:       ctx.Bool("push-target"),
			AutoLabel:        ctx.Bool("auto-label"),
		},
		Extra: kaniko.Extra{
			Executor: ctx.StringSlice("executor-extra-args"),
			Warmer:   ctx.StringSlice("warmer-extra-args"),
		},
	}
}
