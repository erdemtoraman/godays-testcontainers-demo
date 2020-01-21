package wait

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type ForSQL struct {
	UrlFromPort func(port nat.Port) string
	Driver      string
	Port        nat.Port
}

//WaitUntilReady repeatedly tries to run "SELECT 1" query on the given port using sql and driver.
// If the it doesn't succeed until the timeout value which defaults to 10 seconds, it will return an error
func (w ForSQL) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	port, err := target.MappedPort(ctx, w.Port)
	if err != nil {
		return fmt.Errorf("target.MappedPort: %v", err)
	}

	db, err := sql.Open(w.Driver, w.UrlFromPort(port))
	if err != nil {
		return fmt.Errorf("sql.Open: %v", err)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.Tick(time.Millisecond * 100):
			if _, err := db.ExecContext(ctx, "SELECT 1"); err != nil {
				continue
			}
			return nil
		}
	}
}
