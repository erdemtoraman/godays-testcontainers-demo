package unit

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strconv"
	"testing"
)

var (
	_repo *userRepo
	_conn *sqlx.DB
)

func TestMain(m *testing.M) {
	log.Println("Starting postgres container...")
	containerPort := nat.Port("5432/tcp")
	log.Println("a")
	postgres, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				Image:        "postgres",
				ExposedPorts: []string{containerPort.Port()},
				Env: map[string]string{
					"POSTGRES_PASSWORD": "pass",
					"POSTGRES_USER":     "user",
				},
				WaitingFor: wait.ForListeningPort(containerPort),
			},
			Started: true, // auto-start the container
		})
	if err != nil {
		log.Fatal("start:", err)
	}
	mappedPort, err := postgres.MappedPort(context.Background(), containerPort)
	if err != nil {
		log.Fatal("map:", err)
	}
	postgresURLTemplate := "postgres://user:pass@localhost:%s?sslmode=disable"
	postgresURL := fmt.Sprintf(postgresURLTemplate, mappedPort.Port())
	log.Printf("Postgres container started, running at:  %s\n", postgresURL)

	_conn, err = sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Fatal("connect:", err)
	}

	if err := runMigrations(_conn); err != nil {
		log.Fatal("runMigrations:", err)
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
	if _, err := _conn.Exec("TRUNCATE users"); err != nil {
		log.Fatalf("Cannot clear db: %v", err)
	}
}
