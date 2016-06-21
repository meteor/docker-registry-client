package registry

import (
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
)

// ErrDigestNotFound is returned when server does not include digest header in response.
var ErrDigestNotFound = errors.New("Docker-Content-Digest header not found")

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
		return "", err
	}

	// trust server returned digest here / no digest veritifcation
	// TODO: at least, check if schema is V2
	header, ok := resp.Header["Docker-Content-Digest"]
	if !ok || len(header) != 1 {
		return "", ErrDigestNotFound
	}
	return header[0], nil
}
