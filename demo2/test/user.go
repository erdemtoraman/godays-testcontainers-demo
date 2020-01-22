package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
)

type UserServiceConfig struct {
	PostgresURL string
	Port        nat.Port
}

func (s UserServiceConfig) StartContainer(ctx context.Context, networkName string) (internalURL, mappedURL string) {
	dir, _ := os.Getwd()
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			FromDockerfile: tc.FromDockerfile{Context: filepath.Join(dir, "userservice")},
			Networks:       []string{networkName},
			NetworkAliases: map[string][]string{
				networkName: {"user-service"},
			},
			Env:          s.env(),
			ExposedPorts: []string{s.Port.Port()},
			WaitingFor:   wait.ForListeningPort(s.Port),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	mappedPort, err := container.MappedPort(ctx, s.Port)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("http://user-service:%s", s.Port.Port()),
		fmt.Sprintf("http://localhost:%s", mappedPort.Port())
}

func (s UserServiceConfig) env() map[string]string {
	return map[string]string{"POSTGRES_URL": s.PostgresURL, "PORT": s.Port.Port()}
}
