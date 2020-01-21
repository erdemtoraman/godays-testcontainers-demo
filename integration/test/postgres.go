package test

import (
	"context"
	"demo-end2end/wait"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"

	"log"
	"strings"
)

type PostgresConfig struct {
	Password string
	User     string
	DB       string
	Port     string
}

func (p PostgresConfig) env() map[string]string {
	return map[string]string{
		"POSTGRES_PASSWORD": p.Password,
		"POSTGRES_USER":     p.User,
		"POSTGRES_DB":       p.DB,
		"POSTGRES_PORT":     p.Port,
	}
}

func (p PostgresConfig) url(port nat.Port) string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", p.User, p.Password, port.Port(), p.DB)
}

func (p PostgresConfig) StartContainer(ctx context.Context, networkName string) (internalURL, mappedURL string) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{p.Port},
			Env:          p.env(),
			Networks:     []string{networkName},
			NetworkAliases: map[string][]string{
				networkName: {"user-service-postgres"},
			},
			WaitingFor: wait.ForSQL(nat.Port(p.Port), "postgres", p.url),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal("start ", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(p.Port))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(p.url(nat.Port(p.Port)), "@localhost:", fmt.Sprintf("@%s:", "user-service-postgres"), 1), p.url(mappedPort)
}
