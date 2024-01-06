// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package kaniko

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/estesp/manifest-tool/v2/pkg/types"
	"github.com/google/go-containerregistry/pkg/name"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"go.megpoid.dev/drone-kaniko/pkg/manifest"
)

const (
	kanikoExecutor = "/kaniko/executor"
	kanikoWarmer   = "/kaniko/warmer"
)

// Settings for the Plugin.
type Settings struct {
	BuildArgs                   []string
	Cache                       bool
	CacheCopyLayers             bool
	CacheDir                    string
	CacheRepo                   string
	CacheTTL                    time.Duration
	Cleanup                     bool
	CompressedCaching           bool
	Compression                 string
	CompressionLevel            int
	Context                     string
	ContextSubPath              string
	CustomPlatform              string
	Destinations                []string
	DigestFile                  string
	Dockerfile                  string
	Force                       bool
	ForceBuildMetadata          bool
	Git                         string
	IgnorePath                  []string
	IgnoreVarRun                bool
	ImageFsExtractRetry         int
	ImageNameTagWithDigestFile  string
	ImageNameWithDigestFile     string
	Insecure                    bool
	InsecurePull                bool
	InsecureRegistries          []string
	KanikoDir                   string
	Labels                      []string
	LogFormat                   string
	LogTimestamp                bool
	NoPush                      bool
	NoPushCache                 bool
	OCILayoutPath               string
	PushRetry                   int
	RegistryCertificates        []string
	RegistryClientCerts         []string
	RegistryMirror              string
	Reproducible                bool
	SingleSnapshot              bool
	SkipDefaultRegistryFallback bool
	SkipPushPermissionCheck     bool
	SkipTLSVerify               bool
	SkipTLSVerifyPull           bool
	SkipTLSVerifyRegistries     []string
	SkipUnusedStages            bool
	SnapshotMode                string
	TarPath                     string
	Target                      string
	UseNewRun                   bool
	Verbosity                   string
	Auth                        Auth
	Main                        Main
	Manifest                    Manifest
	Extra                       Extra
}

// Auth settings for the Plugin.
type Auth struct {
	Registry string
	Username string
	Password string
	Config   string
}

// Main args for the Plugin.
type Main struct {
	BuildArgsFromEnv    []string
	Debug               bool
	DryRun              bool
	ForceCache          bool
	Tags                []string
	Platforms           []string
	TagsAuto            bool
	TagsSuffix          string
	Images              []string
	Repo                string
	LabelSchema         []string
	Mirror              string
	PushTarget          bool
	AutoLabel           bool
	IncludePlatformTags bool
}

type Manifest struct {
	IgnoreMissing bool
}

// Extra args for the plugin
type Extra struct {
	Executor []string
	Warmer   []string
}

type Platform struct {
	OS           string
	Architecture string
	TarballPath  string
}

func (p *pluginImpl) Validate() error {
	if err := enableCompatibilityMode(&p.settings, &p.pipeline); err != nil {
		return err
	}

	if !p.settings.NoPush && len(p.settings.Destinations) == 0 && p.settings.Main.Repo == "" {
		return errors.New("must provide either no-push or at least one repo/destination")
	}

	for _, entry := range p.settings.RegistryCertificates {
		parts := strings.Split(entry, "=")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid registry-certificate: %s", entry)
		}
	}

	for _, entry := range p.settings.RegistryClientCerts {
		parts := strings.Split(entry, "=")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid registry-client-cert: %s", entry)
		}
	}

	if len(p.settings.Main.Platforms) > 0 {
		if p.settings.CustomPlatform != "" {
			return fmt.Errorf("platforms is set, custom-platform will be ignored")
		}
		if p.settings.TarPath != "" {
			return fmt.Errorf("platforms is set, tar-path will be ignored")
		}
	}

	if err := generateAuthFile(&p.settings.Auth); err != nil {
		return fmt.Errorf("failed to generate docker auth file: %w", err)
	}

	if p.settings.Main.AutoLabel {
		generateLabelSchemas(&p.settings, &p.pipeline)
	}
	// set defaults
	addProxyBuildArgs(&p.settings)
	addArgsFromEnv(&p.settings)

	if p.settings.Cache {
		p.settings.Main.Images = append(p.settings.Main.Images, p.settings.Destinations...)
	}

	if p.settings.Context == "" {
		wd, err := os.Getwd()
		if err != nil {
			p.settings.Context = "."
		} else {
			p.settings.Context = wd
		}
	}

	return nil
}

func (p *pluginImpl) Execute() error {
	var cmds []*exec.Cmd
	cmds = append(cmds, commandKanikoVersion()) // kaniko version

	// no platforms, just build and push directly without a manifest
	if len(p.settings.Main.Platforms) == 0 {
		if len(p.settings.Main.Images) > 0 {
			cmds = append(cmds, commandWarmer(&p.settings)) // kaniko warmer
		}

		cmds = append(cmds, commandBuild(&p.settings)) // kaniko build/push
		if err := runCmds(cmds); err != nil {
			return err
		}

		return nil
	}

	if len(p.settings.Main.Images) > 0 {
		for _, platform := range p.settings.Main.Platforms {
			p.settings.CustomPlatform = platform

			// warmer is called once per platform
			cmds = append(cmds, commandWarmer(&p.settings)) // kaniko warmer
		}
	}

	// list of repositories with their tags
	repositories := make(map[string][]string)

	for idx, destination := range p.settings.Destinations {
		tag, err := name.NewTag(destination)
		if err != nil {
			return fmt.Errorf("invalid destination: %s", destination)
		}

		repoName := tag.Repository.Name()
		if _, ok := repositories[repoName]; !ok {
			repositories[repoName] = []string{tag.TagStr()}
		} else {
			repositories[repoName] = append(repositories[repoName], tag.TagStr())
		}

		// update the destination list with default tags, if missing
		p.settings.Destinations[idx] = tag.Name()
	}

	for _, platform := range p.settings.Main.Platforms {
		p.settings.CustomPlatform = platform

		// kaniko is called once per platform
		cmds = append(cmds, commandBuild(&p.settings)) // kaniko build
	}

	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			var exErr *exec.ExitError
			// ignore warmer errors since the first run there won't be a cache
			if cmd.Args[0] == kanikoWarmer && errors.As(err, &exErr) {
				continue
			}
			return err
		}
	}

	cfg := manifest.Config{
		Insecure: p.settings.Insecure,
	}

	for repoName, tags := range repositories {
		var images []types.ManifestEntry

		for _, platform := range p.settings.Main.Platforms {
			OS, arch, found := strings.Cut(platform, "/")
			if !found {
				return fmt.Errorf("invalid platform: %s", platform)
			}

			images = append(images, types.ManifestEntry{
				Image: repoName + ":" + tags[0] + "-" + arch,
				Platform: ocispec.Platform{
					Architecture: arch,
					OS:           OS,
				},
			})
		}

		target := repoName + ":" + tags[0]

		// push the manifest to the registry, per repository
		manifestErr := manifest.Push(target, tags[1:], images, cfg)
		if manifestErr != nil {
			return fmt.Errorf("failed to push manifest: %w", manifestErr)
		}
	}

	return nil
}

func runCmds(cmds []*exec.Cmd) error {
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			var exerr *exec.ExitError
			// ignore warmer errors since the first run there won't be a cache
			if cmd.Args[0] == kanikoWarmer && errors.As(err, &exerr) {
				continue
			}
			return err
		}
	}

	return nil
}

func commandKanikoVersion() *exec.Cmd {
	return exec.Command(kanikoExecutor, "version")
}

func commandBuild(settings *Settings) *exec.Cmd {
	var args []string
	for _, entry := range settings.BuildArgs {
		args = append(args, "--build-arg", entry)
	}
	if settings.Cache {
		args = append(args, "--cache")
		if settings.CacheRepo != "" {
			args = append(args, "--cache-repo", settings.CacheRepo)
		}
	}
	if settings.CacheCopyLayers {
		args = append(args, "--cache-copy-layers")
	}
	if settings.CacheDir != "" {
		args = append(args, "--cache-dir", settings.CacheDir)
	}
	if settings.CacheTTL != 0 {
		args = append(args, "--cache-ttl", settings.CacheTTL.String())
	}
	if settings.Cleanup {
		args = append(args, "--cleanup")
	}
	if !settings.CompressedCaching {
		args = append(args, "--compressed-caching", "false")
	}
	if settings.Compression != "" {
		args = append(args, "--compression", settings.Compression)
	}
	if settings.CompressionLevel != -1 {
		args = append(args, "--compression-level", strconv.FormatInt(int64(settings.CompressionLevel), 10))
	}
	if settings.Context != "" {
		args = append(args, "--context", settings.Context)
	}
	if settings.ContextSubPath != "" {
		args = append(args, "--context-sub-path", settings.ContextSubPath)
	}
	if settings.CustomPlatform != "" {
		args = append(args, "--custom-platform", settings.CustomPlatform)
	}
	for _, entry := range settings.Destinations {
		if settings.Main.IncludePlatformTags {
			_, arch, _ := strings.Cut(settings.CustomPlatform, "/")
			entry = entry + "-" + arch
		}

		args = append(args, "--destination", entry)
	}
	if settings.DigestFile != "" {
		args = append(args, "--digest-file", settings.DigestFile)
	}
	if settings.Dockerfile != "" {
		args = append(args, "--dockerfile", settings.Dockerfile)
	}
	if settings.Force {
		args = append(args, "--force")
	}
	if settings.ForceBuildMetadata {
		args = append(args, "--force-build-metadata")
	}
	if settings.Git != "" {
		args = append(args, "--git", settings.Git)
	}
	for _, entry := range settings.IgnorePath {
		args = append(args, "--ignore-path", entry)
	}
	if settings.IgnoreVarRun {
		args = append(args, "--ignore-var-run")
	}
	if settings.ImageFsExtractRetry > 0 {
		args = append(args, "--image-fs-extract-retry", strconv.FormatInt(int64(settings.ImageFsExtractRetry), 10))
	}
	if settings.ImageNameTagWithDigestFile != "" {
		args = append(args, "--image-name-tag-with-digest-file", settings.ImageNameTagWithDigestFile)
	}
	if settings.ImageNameWithDigestFile != "" {
		args = append(args, "--image-name-with-digest-file", settings.ImageNameWithDigestFile)
	}
	if settings.Insecure {
		args = append(args, "--insecure")
	}
	if settings.InsecurePull {
		args = append(args, "--insecure-pull")
	}
	for _, entry := range settings.InsecureRegistries {
		args = append(args, "--insecure-registry", entry)
	}
	if settings.KanikoDir != "" {
		args = append(args, "--kaniko-dir", settings.KanikoDir)
	}
	for _, entry := range settings.Labels {
		args = append(args, "--label", entry)
	}
	if settings.LogFormat != "" {
		args = append(args, "--log-format", settings.LogFormat)
	}
	if settings.LogTimestamp {
		args = append(args, "--log-timestamp")
	}
	if settings.NoPush {
		args = append(args, "--no-push")
	}
	if settings.NoPushCache {
		args = append(args, "--no-push-cache")
	}
	if settings.OCILayoutPath != "" {
		args = append(args, "--oci-layout-path", settings.OCILayoutPath)
	}
	if settings.PushRetry > 0 {
		args = append(args, "--push-retry", strconv.FormatInt(int64(settings.PushRetry), 10))
	}
	for _, entry := range settings.RegistryCertificates {
		args = append(args, "--registry-certificate", entry)
	}
	for _, entry := range settings.RegistryClientCerts {
		args = append(args, "--registry-client-cert", entry)
	}
	if settings.RegistryMirror != "" {
		args = append(args, "--registry-mirror", settings.RegistryMirror)
	}
	if settings.Reproducible {
		args = append(args, "--reproducible")
	}
	if settings.SingleSnapshot {
		args = append(args, "--single-snapshot")
	}
	if settings.SkipDefaultRegistryFallback {
		args = append(args, "--skip-default-registry-fallback")
	}
	if settings.SkipPushPermissionCheck {
		args = append(args, "--skip-push-permission-check")
	}
	if settings.SkipTLSVerify {
		args = append(args, "--skip-tls-verify")
	}
	if settings.SkipTLSVerifyPull {
		args = append(args, "--skip-tls-verify-pull")
	}
	for _, entry := range settings.SkipTLSVerifyRegistries {
		args = append(args, "--skip-tls-verify-registry", entry)
	}
	if settings.SkipUnusedStages {
		args = append(args, "--skip-unused-stages")
	}
	if settings.SnapshotMode != "" {
		args = append(args, "--snapshotMode", settings.SnapshotMode)
	}
	if settings.TarPath != "" {
		args = append(args, "--tar-path", settings.TarPath)
	}
	if settings.Target != "" {
		args = append(args, "--target", settings.Target)
	}
	if settings.UseNewRun {
		args = append(args, "--use-new-run")
	}
	if settings.Verbosity != "" {
		args = append(args, "--verbosity", settings.Verbosity)
	}
	if len(settings.Extra.Executor) > 0 {
		args = append(args, settings.Extra.Executor...)
	}
	return exec.Command(kanikoExecutor, args...)
}

func commandWarmer(settings *Settings) *exec.Cmd {
	var args []string
	if settings.CacheDir != "" {
		args = append(args, "--cache-dir", settings.CacheDir)
	}
	if settings.CacheTTL != 0 {
		args = append(args, "--cache-ttl", settings.CacheTTL.String())
	}
	if settings.CustomPlatform != "" {
		args = append(args, "--customPlatform", settings.CustomPlatform)
	}
	if settings.Main.ForceCache {
		args = append(args, "--force")
	}
	for _, entry := range settings.Main.Images {
		args = append(args, "--image", entry)
	}
	if settings.InsecurePull {
		args = append(args, "--insecure-pull")
	}
	for _, entry := range settings.InsecureRegistries {
		args = append(args, "--insecure-registry", entry)
	}
	if settings.LogFormat != "" {
		args = append(args, "--log-format", settings.LogFormat)
	}
	if settings.LogTimestamp {
		args = append(args, "--log-timestamp")
	}
	for _, entry := range settings.RegistryCertificates {
		args = append(args, "--registry-certificate", entry)
	}
	if settings.RegistryMirror != "" {
		args = append(args, "--registry-mirror", settings.RegistryMirror)
	}
	if settings.SkipTLSVerifyPull {
		args = append(args, "--skip-tls-verify-pull")
	}
	for _, entry := range settings.SkipTLSVerifyRegistries {
		args = append(args, "--skip-tls-verify-registry", entry)
	}
	if settings.Verbosity != "" {
		args = append(args, "--verbosity", settings.Verbosity)
	}
	if len(settings.Extra.Warmer) > 0 {
		args = append(args, settings.Extra.Warmer...)
	}
	return exec.Command(kanikoWarmer, args...)
}
