package registry

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
)

// ErrBadSchemaVersion is returned when server does not support schema V2.
var ErrBadSchemaVersion = errors.New("Server does not seem to support schema version 2")

// ErrDigestNotFound is returned when server does not include digest header in response.
var ErrDigestNotFound = errors.New("Docker-Content-Digest header not found")

type manifestReponse struct {
	SchemaVersion int `json:"schemaVersion"`
}

func (registry *Registry) GetManifestDigest(ctx context.Context, repository string, reference string) (string, error) {
	url := registry.url("/v2/%s/manifests/%s", repository, reference)
	registry.Logf("registry.manifests url=%s repository=%s reference=%s", url, repository, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header["Accept"] = []string{"application/vnd.docker.distribution.manifest.v2+json"}

	resp, err := ctxhttp.Do(ctx, registry.Client, req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", chooseError(ctx, err)
	}

	var manifest manifestReponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&manifest)
	if err != nil {
		return "", err
	}

	if manifest.SchemaVersion != 2 {
		return "", ErrBadSchemaVersion
	}

	header, ok := resp.Header["Docker-Content-Digest"]
	if !ok || len(header) != 1 {
		return "", ErrDigestNotFound
	}

	// No verification is done to check if this digest actually match
	// the content of returned image manifest. This trusts the server
	// to return the correct value.
	return header[0], nil
}
