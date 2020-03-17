package kaniko

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const kanikoExecutor = "/kaniko/executor"
const kanikoWarmer = "/kaniko/warmer"

// Settings for the Plugin.
type Settings struct {
	BuildArgs               []string
	Cache                   bool
	CacheDir                string
	CacheRepo               string
	CacheTTL                time.Duration
	Cleanup                 bool
	Context                 string
	Destinations            []string
	DigestFile              string
	Dockerfile              string
	Force                   bool
	ImageNameWithDigestFile string
	Insecure                bool
	InsecurePull            bool
	InsecureRegistries      []string
	Labels                  []string
	LogFormat               string
	NoPush                  bool
	OCILayoutPath           string
	RegistryCertificates    []string
	RegistryMirror          string
	Reproducible            bool
	SingleSnapshot          bool
	SkipTLSVerify           bool
	SkipTLSVerifyPull       bool
	SkipTLSVerifyRegistries []string
	SnapshotMode            string
	TarPath                 string
	Target                  string
	Verbosity               string
	WhitelistVarRun         bool
	// other args
	BuildArgsFromEnv []string
}

func (p *pluginImpl) Validate() error {
	if p.settings.NoPush && len(p.settings.Destinations) == 0 {
		return errors.New("must provide either no-push or at least one destination")
	}
	for _, entry := range p.settings.RegistryCertificates {
		parts := strings.Split(entry, "=")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid registry-certificate: %s", entry)
		}
	}
	return nil
}

func (p *pluginImpl) Execute() error {
	addProxyBuildArgs(&p.settings)

	var cmds []*exec.Cmd
	cmds = append(cmds, commandVersion())          // kaniko version
	cmds = append(cmds, commandBuild(&p.settings)) // kaniko build

	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func commandVersion() *exec.Cmd {
	return exec.Command(kanikoExecutor, "version")
}

func commandBuild(settings *Settings) *exec.Cmd {
	var args []string

	for _, entry := range settings.BuildArgsFromEnv {
		addProxyValue(settings, entry)
	}
	for _, entry := range settings.BuildArgs {
		args = append(args, "--build-arg", entry)
	}
	if settings.Cache {
		args = append(args, "--cache", "true")
	}
	if settings.CacheDir != "" {
		args = append(args, "--cache-dir", settings.CacheDir)
	}
	if settings.CacheRepo != "" {
		args = append(args, "--cache-repo", settings.CacheRepo)
	}
	if settings.CacheTTL != 0 {
		args = append(args, "--cache-ttl", settings.CacheTTL.String())
	}
	if settings.Cleanup {
		args = append(args, "--cleanup")
	}
	if settings.Context != "" {
		args = append(args, "--context", settings.Context)
	}
	for _, entry := range settings.Destinations {
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
	for _, entry := range settings.Labels {
		args = append(args, "--label", entry)
	}
	if settings.LogFormat != "" {
		args = append(args, "--log-format", settings.LogFormat)
	}
	if settings.NoPush {
		args = append(args, "--no-push")
	}
	if settings.OCILayoutPath != "" {
		args = append(args, "--oci-layout-path", settings.OCILayoutPath)
	}
	for _, entry := range settings.RegistryCertificates {
		args = append(args, "--registry-certificate", entry)
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
	if settings.SkipTLSVerify {
		args = append(args, "--skip-tls-verify")
	}
	if settings.SkipTLSVerifyPull {
		args = append(args, "--skip-tls-verify-pull")
	}
	for _, entry := range settings.SkipTLSVerifyRegistries {
		args = append(args, "--skip-tls-verify-registry", entry)
	}
	if settings.SnapshotMode != "" {
		args = append(args, "--snapshot-mode", settings.SnapshotMode)
	}
	if settings.TarPath != "" {
		args = append(args, "--tar-path", settings.TarPath)
	}
	if settings.Target != "" {
		args = append(args, "--target", settings.Target)
	}
	if settings.Verbosity != "" {
		args = append(args, "--verbosity", settings.Verbosity)
	}
	if settings.WhitelistVarRun {
		args = append(args, "--whitelist-var-run")
	}

	return exec.Command(kanikoExecutor, args...)
}

func trace(cmd *exec.Cmd) {
	_, _ = fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

// helper function to add proxy values from the environment
func addProxyBuildArgs(settings *Settings) {
	addProxyValue(settings, "http_proxy")
	addProxyValue(settings, "https_proxy")
	addProxyValue(settings, "no_proxy")
}

// helper function to add the upper and lower case version of a proxy value.
func addProxyValue(settings *Settings, key string) {
	value := getProxyValue(key)

	if len(value) > 0 && !hasProxyBuildArg(settings, key) {
		settings.BuildArgs = append(settings.BuildArgs, fmt.Sprintf("%s=%s", key, value))
		settings.BuildArgs = append(settings.BuildArgs, fmt.Sprintf("%s=%s", strings.ToUpper(key), value))
	}
}

// helper function to get a proxy value from the environment.
//
// assumes that the upper and lower case versions of are the same.
func getProxyValue(key string) string {
	value := os.Getenv(key)

	if len(value) > 0 {
		return value
	}

	return os.Getenv(strings.ToUpper(key))
}

// helper function that looks to see if a proxy value was set in the build args.
func hasProxyBuildArg(settings *Settings, key string) bool {
	keyUpper := strings.ToUpper(key)

	for _, s := range settings.BuildArgs {
		if strings.HasPrefix(s, key) || strings.HasPrefix(s, keyUpper) {
			return true
		}
	}

	return false
}
