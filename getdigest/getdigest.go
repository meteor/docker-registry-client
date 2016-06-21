package main

import (
	"flag"
	"fmt"
	"github.com/meteor/docker-registry-client/registry"
	"golang.org/x/net/context"
	"os"
	"strings"
)

var (
	username = flag.String("username", "", "Registry username")
	password = flag.String("password", "", "Registry password")
	hostname = flag.String("registry", "registry-1.docker.io", "Registry hostname")
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <repository:tag>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	split := strings.Split(flag.Arg(0), ":")
	if len(split) != 2 {
		fmt.Fprintf(os.Stderr, "Need a docker image spec in form of <repo>:<tag>\n")
		os.Exit(1)
	}

	url := fmt.Sprintf("https://%s/", *hostname)
	hub, err := registry.New(context.Background(), url, *username, *password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	digest, err := hub.GetManifestDigest(context.Background(), split[0], split[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(digest)
}
