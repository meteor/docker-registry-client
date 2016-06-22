package registry

import (
	"golang.org/x/net/context"
	"testing"
)

func buildAnonymousDockerHub(t *testing.T, ctx context.Context) *Registry {
	r, err := New(ctx, "https://registry-1.docker.io/", "", "")
	if err != nil {
		t.Error(err)
	}
	return r
}

func TestTags(t *testing.T) {
	r := buildAnonymousDockerHub(t, context.Background())

	tags, err := r.Tags(context.Background(), "library/ubuntu")
	if err != nil {
		t.Error(err)
	}

	tagsMap := make(map[string]struct{})
	for _, tag := range tags {
		tagsMap[tag] = struct{}{}
	}

	if _, ok := tagsMap["latest"]; !ok {
		t.Fail()
	}
	if _, ok := tagsMap["16.04"]; !ok {
		t.Fail()
	}
}

func TestManifestDigest(t *testing.T) {
	r := buildAnonymousDockerHub(t, context.Background())

	// if we ask for a manifest with this specific digest, it should
	// return the exact same digest
	expected := "sha256:4a731fb46adc5cefe3ae374a8b6020fc1b6ad667a279647766e9a3cd89f6fa92"
	digest, err := r.GetManifestDigest(context.Background(), "library/busybox", expected)
	if err != nil {
		t.Error(err)
	}

	if expected != digest {
		t.Fail()
	}
}
