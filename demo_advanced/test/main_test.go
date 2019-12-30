package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

const networkName = "integration-test-network"

var userServiceURL, ticketServiceURL string

func init() {
	provider, err := testcontainers.ProviderDocker.GetProvider()
	if err != nil {
		log.Fatal(err)
	}

	if _, err = provider.GetNetwork(context.Background(), testcontainers.NetworkRequest{Name: networkName, Driver: "bridge"}); err != nil {
		if _, err := provider.CreateNetwork(context.Background(), testcontainers.NetworkRequest{Name: networkName, Driver: "bridge"}); err != nil {
			log.Fatal(err)
		}
	}
}

func TestMain(m *testing.M) {
	os.Chdir("..")

	postgresConfig := PostgresConfig{
		Password: "password",
		User:     "postgres",
		DB:       "userservice",
		Port:     "5432",
	}

	postgresInternal, mappedPostgres := postgresConfig.Start(context.Background(), networkName)
	log.Println("postgres running at: ", mappedPostgres)
	internalUser, mappedUser := UserServiceConfig{PostgresURL: postgresInternal, Port: "8080"}.StartDocker(context.Background(), networkName)

	log.Println("user service running at: ", mappedUser)

	_, ticketServiceURL = TicketServiceConfig{UserServiceURL: internalUser, Port: "8080"}.StartDocker(context.Background(), networkName)

	log.Println("ticket service running at: ", mappedUser)

	userServiceURL = mappedUser

	os.Exit(m.Run())
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TicketPost struct {
	UserID int64  `json:"user_id"`
	Movie  string `json:"movie"`
}
type Ticket struct {
	ID    string `json:"id"`
	Movie string `json:"movie"`
	User  User   `json:"user"`
}

func Test_Integrations(t *testing.T) {
	var createdUser User
	t.Run("create and get a user", func(t *testing.T) {
		resp, _ := http.Post(userServiceURL+"/users", "application/json", structToJsonReader(User{Name: "berliner"}))
		require.NoError(t, responseToStruct(resp.Body, &createdUser))
		assert.Equal(t, "berliner", createdUser.Name)

		var getUser User
		resp, _ = http.Get(fmt.Sprintf("%s/users/%d", userServiceURL, createdUser.ID))
		require.NoError(t, responseToStruct(resp.Body, &getUser))
		assert.Equal(t, createdUser, getUser)
	})

	t.Run("create a ticket", func(t *testing.T) {
		resp, _ := http.Post(ticketServiceURL+"/tickets", "application/json", structToJsonReader(TicketPost{Movie: "dogs of berlin", UserID: createdUser.ID}))
		var ticket struct {
			ID string `json:"id"`
		}
		require.NoError(t, responseToStruct(resp.Body, &ticket))
		assert.Equal(t, "berliner", createdUser.Name)

		var getTicket Ticket
		resp, _ = http.Get(fmt.Sprintf("%s/tickets/%s", ticketServiceURL, ticket.ID))
		require.NoError(t, responseToStruct(resp.Body, &getTicket))
		assert.Equal(t, createdUser, getTicket.User)
		assert.Equal(t, ticket.ID, getTicket.ID)
		assert.Equal(t, "dogs of berlin", getTicket.Movie)
	})
}

func structToJsonReader(v interface{}) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func responseToStruct(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}
