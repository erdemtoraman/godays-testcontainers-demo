package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type TicketServiceConfig struct {
	UserServiceURL string
	Port           nat.Port
}

func (t TicketServiceConfig) StartContainer(ctx context.Context, networkName string) (internalURL, mappedURL string) {
	dir, _ := os.Getwd()
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			FromDockerfile: tc.FromDockerfile{Context: filepath.Join(dir, "ticketservice")},
			Networks:       []string{networkName},
			NetworkAliases: map[string][]string{networkName: {"ticket-service"}},
			Env:            t.env(),
			ExposedPorts:   []string{t.Port.Port()},
			WaitingFor: wait.
				ForHTTP("/health").
				WithPort(t.Port).
				WithStatusCodeMatcher(func(status int) bool {
					return status == http.StatusOK
				}),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	mappedPort, err := container.MappedPort(ctx, t.Port)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("http://ticket-service:%s", t.Port),
		fmt.Sprintf("http://localhost:%s", mappedPort.Port())
}

func (t TicketServiceConfig) env() map[string]string {
	return map[string]string{"USER_SERVICE_URL": t.UserServiceURL, "PORT": t.Port.Port()}
}
