package unit

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strconv"
	"testing"
)

var (
	_ctx  = context.Background()
	_repo UserRepo
	_conn *sqlx.DB

	containerPort    nat.Port = "5432"
	containerRequest          = testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres",
			ExposedPorts: []string{containerPort.Port()},
			Env:          map[string]string{"POSTGRES_PASSWORD": "pass", "POSTGRES_USER": "user", "POSTGRES_DB": "ourservice"},
			WaitingFor:   wait.ForListeningPort("5432"),
		},
		Started: true,
	}
)

func TestMain(m *testing.M) {
	container, err := testcontainers.GenericContainer(_ctx, containerRequest)
	if err != nil {
		log.Fatal("start:", err)
	}
	mappedPort, err := container.MappedPort(_ctx, containerPort)
	if err != nil {
		log.Fatal("map:", err)
	}

	connection, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://user:pass@localhost:%s/ourservice?sslmode=disable", mappedPort.Port()))
	if err != nil {
		log.Fatal("connect:", err)
	}
	_conn = connection

	if err := runMigrations(_conn); err != nil {
		log.Fatal("runMigrations", err)
	}

	_repo = NewRepo(_conn)
	os.Exit(m.Run())
}

func TestRepoImp_CreateUser(t *testing.T) {
	truncateDB()

	user, err := _repo.CreateUser("username")
	require.NoError(t, err)
	assert.Equal(t, "username", user.Name)
	assert.NotZero(t, user.ID)

	user, err = _repo.CreateUser("username")
	assert.Error(t, err, "names are unique")

}

func TestRepoImp_GetUsers(t *testing.T) {
	t.Run("get single user", func(t *testing.T) {
		truncateDB()
		user, err := _repo.CreateUser("username")
		require.NoError(t, err)

		getUser, err := _repo.GetUserByID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, user, getUser)
	})

	t.Run("get all users", func(t *testing.T) {
		truncateDB()

		for i := 0; i < 10; i++ {
			_, err := _repo.CreateUser(strconv.Itoa(i))
			require.NoError(t, err)
		}
		users, err := _repo.GetAllUsers()
		require.NoError(t, err)
		assert.Len(t, users, 10)
	})

}

//noinspection SqlResolve
func truncateDB() {
	_, err := _conn.Exec("TRUNCATE users")
	if err != nil {
		log.Fatalf("Cannot clear db: %v", err)
	}
}
