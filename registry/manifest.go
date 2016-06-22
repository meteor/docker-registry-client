package registry

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
)

// ErrBadSchema is returned when server returns an unsupported or invalid schema format.
var ErrBadSchema = errors.New("Server did not return a supported schema format")

// ErrDigestNotFound is returned when server does not include digest header in response.
var ErrDigestNotFound = errors.New("Docker-Content-Digest header not found")

type manifestResponse struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"` // v2 only
}

const (
	schemaV1MimeType = "application/vnd.docker.distribution.manifest.v1+json"
	schemaV2MimeType = "application/vnd.docker.distribution.manifest.v2+json"
)

func (registry *Registry) GetManifestDigest(ctx context.Context, repository string, reference string) (string, error) {
	url := registry.url("/v2/%s/manifests/%s", repository, reference)
	registry.Logf("registry.manifests url=%s repository=%s reference=%s", url, repository, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// both schema formats are valid
	req.Header["Accept"] = []string{schemaV2MimeType, schemaV1MimeType}

	resp, err := ctxhttp.Do(ctx, registry.Client, req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", chooseError(ctx, err)
	}

	var manifest manifestResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&manifest); err != nil {
		return "", err
	}

	// sanity check to at least make sure we are getting an image manifest
	if manifest.SchemaVersion == 1 {
		// no check for schema version 1
	} else if manifest.SchemaVersion == 2 {
		// In version 2, server may return a list of architecture-dependent
		// manifests. This makes sure we're only getting the default (linux,
		// arm64) image manifest.
		if manifest.MediaType != schemaV2MimeType {
			return "", ErrBadSchema
		}
	} else {
		return "", ErrBadSchema
	}

	header, ok := resp.Header["Docker-Content-Digest"]
	if !ok || len(header) != 1 {
		return "", ErrDigestNotFound
	}

	// No verification is done to check if this digest actually match
	// the content of returned image manifest. This implicitly trusts
	// registry to return correct value.
	return header[0], nil
}
