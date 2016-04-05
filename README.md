# Fork of http://github.com/heroku/docker-registry-client/ with fewer features and dependencies

http://github.com/heroku/docker-registry-client/registry is a client for the V2
Docker Registry API. It currently depends on an old version of
http://github.com/docker/distribution which makes it difficult to build.

For our project's needs, we don't actually need the parts of the upstream client
that depend on Docker's library.  So
`github.com/meteor/docker-registry-client/registry` lacks the Layer and Manifest
APIs of the upstream library; it only contains the Tags APIs. But it doesn't
require you to vendor in an old version of http://github.com/docker/distribution
using pre-core vendoring support!

# Docker Registry Client

An API client for the [V2 Docker Registry
API](http://docs.docker.com/registry/spec/api/), for Go applications.

## Imports

```go
import (
    "github.com/meteor/docker-registry-client/registry"
)
```

## Creating A Client

```go
url      := "https://registry-1.docker.io/"
username := "" // anonymous
password := "" // anonymous
hub, err := registry.New(url, username, password)
```

Creating a registry will also ping it to verify that it supports the registry
API, which may fail. Failures return non-`nil` err values.

Authentication supports both HTTP Basic authentication and OAuth2 token
negotiation.

## Listing Tags

Each Docker repository has a set of tags -- named images that can be downloaded.

```go
tags, err := hub.Tags("heroku/cedar")
```

The tags will be returned as a slice of `string`s.
