package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	waitG "github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
)

type TicketServiceConfig struct {
	UserServiceURL string
	Port           string
}

func (t TicketServiceConfig) StartDocker(ctx context.Context, networkName string) (internalURL, mappedURL string) {
	dir, _ := os.Getwd()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{Context: filepath.Join(dir, "ticketservice")},
			Networks:       []string{networkName},
			NetworkAliases: map[string][]string{networkName: {"ticket-service"}},
			Env:            t.env(),
			ExposedPorts:   []string{t.Port},
			WaitingFor:     waitG.ForListeningPort(nat.Port(t.Port)),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	mappedPort, err := container.MappedPort(ctx, nat.Port(t.Port))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("http://%s:%s", "ticket-service", t.Port), fmt.Sprintf("http://localhost:%s", mappedPort.Port())
}

func (t TicketServiceConfig) env() map[string]string {
	return map[string]string{"USER_SERVICE_URL": t.UserServiceURL, "PORT": t.Port}
}
