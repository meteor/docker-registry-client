package registry

import (
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

func (registry *Registry) getJson(ctx context.Context, url string, response interface{}) error {
	resp, err := ctxhttp.Get(ctx, registry.Client, url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return chooseError(ctx, err)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(response); err != nil {
		return err
	}

	return nil
}
