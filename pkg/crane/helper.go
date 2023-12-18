package crane

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

type config struct {
	UseDigest bool
}

type Option func(settings *config)

func WithDigest() Option {
	return func(settings *config) {
		settings.UseDigest = true
	}
}

func Push(file string, opts ...Option) (string, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	manifest, err := tarball.LoadManifest(pathOpener(file))
	if err != nil {
		return "", fmt.Errorf("failed to load manifest: %w", err)
	}

	if len(manifest) != 1 {
		return "", fmt.Errorf("manifest must contain only one image")
	}

	descriptor := manifest[0]
	var digest string
	repos := make(map[string]string)

	for _, t := range descriptor.RepoTags {
		tag, err := name.NewTag(t)
		if err != nil {
			return "", fmt.Errorf("failed to parse tag: %w", err)
		}

		img, err := tarball.Image(pathOpener(file), &tag)
		if err != nil {
			return "", fmt.Errorf("failed to load image: %w", err)
		}

		hash, err := img.Digest()
		if err != nil {
			return "", fmt.Errorf("failed to get digest: %w", err)
		}

		// only return the digest of the first tag since they all point to the same image
		if digest == "" {
			digest = hash.String()
		}

		repoName := tag.Repository.Name()
		if _, ok := repos[repoName]; !ok {
			repos[repoName] = hash.String()
		} else if repos[repoName] == hash.String() {
			// skip if the tag points to the same image
			continue
		}

		var target string

		if cfg.UseDigest {
			target = fmt.Sprintf(repoName + "@" + digest)
		} else {
			target = fmt.Sprintf(repoName + ":" + tag.TagStr())
		}

		if err = crane.Push(img, target); err != nil {
			return "", fmt.Errorf("failed to push image %s: %w", tag.String(), err)
		}

		slog.Info("Image pushed to registry", "target", target)
	}

	if digest == "" {
		return "", fmt.Errorf("there are no tagged images in the manifest")
	}

	return digest, nil
}

func pathOpener(path string) tarball.Opener {
	return func() (io.ReadCloser, error) {
		return os.Open(path)
	}
}
