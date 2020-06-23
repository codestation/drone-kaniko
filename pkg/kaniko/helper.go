package kaniko

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	tags "github.com/drone-plugins/drone-docker"
	"github.com/drone-plugins/drone-plugin-lib/drone"
	"github.com/sirupsen/logrus"
)

type authConfig struct {
	Auths map[string]authEntry `json:"auths"`
}

type authEntry struct {
	Auth string `json:"auth"`
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

func generateAuthFile(settings *Settings) error {
	if settings.Username != "" && settings.Password != "" {
		encodedPassword := base64.StdEncoding.EncodeToString([]byte(settings.Username + ":" + settings.Password))

		auth := authConfig{
			Auths: map[string]authEntry{
				settings.Registry: {Auth: encodedPassword},
			},
		}

		data, err := json.MarshalIndent(auth, "", "\t")
		if err != nil {
			return err
		}

		configJson := "/kaniko/.docker/config.json"
		logrus.Infof("Generating auth info in %s", configJson)
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

func enableCompatibilityMode(settings *Settings, pipeline *drone.Pipeline) error {
	if settings.Insecure {
		settings.InsecurePull = true
		settings.InsecurePush = true
	}
	if settings.Debug {
		settings.Verbosity = "debug"
	}

	if settings.Mirror != "" {
		if settings.RegistryMirror != "" {
			return errors.New("mirror and registry-mirror cannot be set at the same time")
		}
		re := regexp.MustCompile("^(http?)://(.*)")
		matches := re.FindStringSubmatch(settings.Mirror)
		if len(matches) > 2 {
			if matches[1] == "http" {
				// mark as insecure
				settings.InsecurePull = true
			}
			// remove scheme from url
			settings.RegistryMirror = matches[2]
		} else {
			return fmt.Errorf("invalid mirror: %s", settings.Mirror)
		}
	}

	if settings.TagsAuto {
		if tags.UseDefaultTag(pipeline.Commit.Ref, pipeline.Repo.Branch) {
			tag, err := tags.DefaultTagSuffix(pipeline.Commit.Ref, settings.TagsSuffix)
			if err != nil {
				logrus.Printf("cannot build docker image for %s, invalid semantic version", pipeline.Commit.Ref)
				return err
			}
			settings.Tags = tag
		} else {
			logrus.Warnf("Skipping automated build for %s", pipeline.Commit.Ref)
		}
	}

	if settings.Repo != "" {
		for _, entry := range settings.Tags {
			dest := fmt.Sprintf("%s:%s", settings.Repo, entry)
			settings.Destinations = append(settings.Destinations, dest)
		}
	}

	return nil
}

func generateLabelSchemas(settings *Settings, pipeline *drone.Pipeline) {
	labelSchema := []string{
		"schema-version=1.0",
		fmt.Sprintf("build-date=%s", time.Now().Format(time.RFC3339)),
		fmt.Sprintf("vcs-ref=%s", pipeline.Commit.SHA),
		fmt.Sprintf("vcs-url=%s", pipeline.Repo.HTTPURL),
	}

	if len(settings.LabelSchema) > 0 {
		labelSchema = append(labelSchema, settings.LabelSchema...)
	}

	for _, label := range labelSchema {
		settings.Labels = append(settings.Labels, fmt.Sprintf("org.label-schema.%s", label))
	}
}
