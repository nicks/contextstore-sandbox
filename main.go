package main

import (
	"log"
	"runtime"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/store"
	"github.com/docker/cli/cli/flags"
	"github.com/pkg/errors"
)

func main() {
	go func() {
		err := ensureDesktopContextExists()
		if err != nil {
			log.Printf("WARNING: %v", err)
		}
		log.Println("finished init")
	}()
	time.Sleep(time.Millisecond / 100)
	log.Fatal("EXITING")
}

func ensureDesktopContextExists() error {
	name := "desktop-linux"
	host := "npipe:////./pipe/dockerDesktopLinuxEngine"
	if runtime.GOOS != "windows" {
		host = "unix://var/run/docker.sock"
	}

	cli, err := initCLI()
	if err != nil {
		return err
	}
	s := cli.ContextStore()
	if _, err := s.GetMetadata(name); err != nil && !store.IsErrContextDoesNotExist(err) {
		return errors.Wrap(err, "querying "+name+" docker context")
	}

	m := store.Metadata{
		Name:     name,
		Metadata: map[string]interface{}{}, // required by docker-compose
		Endpoints: map[string]interface{}{
			"docker": map[string]interface{}{
				"Host":          host,
				"SkipTLSVerify": false,
			},
		},
	}
	if err := s.CreateOrUpdate(m); err != nil {
		return errors.Wrap(err, "creating "+name+" context")
	}
	log.Printf("create docker cli context %s", name)
	return nil
}

func initCLI() (*command.DockerCli, error) {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, errors.Wrap(err, "creating the Docker CLI")
	}
	if err := dockerCli.Initialize(&flags.ClientOptions{}); err != nil {
		return nil, errors.Wrap(err, "initializing the Docker CLI")
	}
	return dockerCli, nil
}
