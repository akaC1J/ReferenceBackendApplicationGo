//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
)

var db *pgx.Conn

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForSQL(
			"5432/tcp",
			"postgres",
			func(host string, port nat.Port) string {
				return fmt.Sprintf("host=%s port=%s user=testuser password=testpassword dbname=testdb sslmode=disable", host, port.Port())
			},
		),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}
	mappedPort, _ := postgresContainer.MappedPort(ctx, "5432/tcp")
	log.Printf("Container started on port: %s", mappedPort.Port())

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpassword dbname=testdb sslmode=disable", host, port.Port())

	db, err = pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("Database connected successfully!")
	os.Setenv("DATABASE_MASTER_HOST_PORT_0", fmt.Sprintf("localhost:%s", port.Port()))
	code := m.Run()

	if err := postgresContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %v", err)
	}

	os.Exit(code)
}
func executeSQLFile(db *pgx.Conn, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %v", filename, err)
	}
	_, err = db.Exec(context.Background(), string(data))
	if err != nil {
		return fmt.Errorf("failed to execute SQL file %s: %v", filename, err)
	}
	return nil
}

func setupTest(t *testing.T) {
	err := executeSQLFile(db, "up.sql")
	if err != nil {
		t.Fatalf("failed to execute up.sql: %v", err)
	}

	t.Cleanup(func() {
		err := executeSQLFile(db, "down.sql")
		if err != nil {
			t.Fatalf("failed to execute down.sql during cleanup: %v", err)
		}
	})
}
