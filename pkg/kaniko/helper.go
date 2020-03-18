package kaniko

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	tags "github.com/drone-plugins/drone-docker"
	"github.com/drone-plugins/drone-plugin-lib/drone"
)

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

func generateAuthFile(settings *Settings) error {
	if settings.Username != "" && settings.Password != "" {
		encodedPassword := base64.StdEncoding.EncodeToString([]byte(settings.Username + ":" + settings.Password))
		var registry string
		if settings.Registry != "" {
			registry = settings.Registry
		} else {
			registry = "https://index.docker.io/v1/"
		}

		auth := authConfig{
			Auths: map[string]authEntry{
				registry: {Auth: encodedPassword},
			},
		}

		data, err := json.MarshalIndent(auth, "", "\t")
		if err != nil {
			return err
		}

		configJson := "/kaniko/.docker/config.json"
		log.Printf("Generating auth info in %s", configJson)
		err = ioutil.WriteFile(configJson, data, 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func addArgsFromEnv(settings *Settings) {
	for _, entry := range settings.BuildArgsFromEnv {
		addProxyValue(settings, entry)
	}
}

func enableCompatibilityMode(settings *Settings, pipeline *drone.Pipeline) bool {
	if settings.Insecure {
		settings.InsecurePull = true
		settings.InsecurePush = true
	}
	if settings.Debug {
		settings.Verbosity = "debug"
	}

	if settings.TagsAuto {
		if tags.UseDefaultTag(pipeline.Commit.Ref, pipeline.Repo.Branch) {
			settings.Tags = tags.DefaultTagSuffix(pipeline.Commit.Ref, settings.TagsSuffix)
		} else {
			log.Printf("Skipping automated build for %s", pipeline.Commit.Ref)
			return false
		}
	}
	if settings.Repo != "" {
		for _, entry := range settings.Tags {
			dest := fmt.Sprintf("%s:%s", settings.Repo, entry)
			log.Printf("Using --destination %s", dest)
			settings.Destinations = append(settings.Destinations)
		}
	}

	return true
}
