package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"api5back/ent"
	"api5back/ent/migrate"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestCase struct {
	Name string
	Run  func(t *testing.T)
}

type IntegrationEnvironment struct {
	Container   testcontainers.Container
	Credentials *Credentials
	Client      *ent.Client
	Error       error
	migrated    bool
}

func startTestingDatabaseContainer(
	ctx context.Context,
	credentials *Credentials,
) (testcontainers.Container, error) {
	var databaseName string
	if credentials.Name != nil {
		databaseName = fmt.Sprintf("%s", *credentials.Name)
	} else {
		databaseName = ""
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		Name:         "khali-api5-TI-postgres",
		ExposedPorts: []string{"5432/tcp"},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"5432/tcp": []nat.PortBinding{{
					HostIP:   "localhost",
					HostPort: fmt.Sprintf("%d/tcp", *credentials.Port),
				}},
			}
		},
		// Wait for _this string_ to appear in the container logs.
		// It's the last unique string that appears, but not the last line.
		// The last line appears the first time way sooner so wouldn't work.
		// I've also tried this, no bueno:
		//		WaitingFor: wait.ForListeningPort("5432/tcp"),
		WaitingFor: wait.ForLog("listening on IPv6 address"),
		Env: map[string]string{
			"POSTGRES_USER":     credentials.User,
			"POSTGRES_PASSWORD": credentials.Pass,
			"POSTGRES_DB":       databaseName,
		},
	}

	// Start the container
	return testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Reuse:            true,
		},
	)
}

func newIntegrationEnvironment(
	ctx context.Context,
	credentials *Credentials,
) *IntegrationEnvironment {
	fmt.Println("integration_environment::newIntegrationEnvironment -- Creating new Integration environment")

	container, err := startTestingDatabaseContainer(ctx, credentials)
	if err != nil {
		fmt.Println("integration_environment::newIntegrationEnvironment -- Error: failed to start the testing database container")

		return &IntegrationEnvironment{
			Error: fmt.Errorf("failed to start the testing database container: %v", err),
		}
	}

	return &IntegrationEnvironment{
		Container:   container,
		Credentials: credentials,
	}
}

func DefaultIntegrationEnvironment(ctx context.Context) *IntegrationEnvironment {
	credentials, err := newTestingCredentials()
	if err != nil {
		return &IntegrationEnvironment{
			Error: fmt.Errorf("failed to get testing database credentials: %v", err),
		}
	}

	return newIntegrationEnvironment(ctx, credentials).
		WithClient()
}

func (intEnv *IntegrationEnvironment) WithSleep() *IntegrationEnvironment {
	if intEnv.Error != nil || intEnv.Container == nil {
		return intEnv
	}

	sleepTimeMs := getContainerConnectionDelayMs()

	fmt.Printf("integration_environment::WithSleep -- Sleeping for %dms to allow the container to connect\n", sleepTimeMs)
	time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)
	fmt.Println("integration_environment::WithSleep -- Done sleeping")

	return intEnv
}

func (intEnv *IntegrationEnvironment) WithClient() *IntegrationEnvironment {
	if intEnv.Error != nil || intEnv.Container == nil || intEnv.Client != nil {
		return intEnv
	}

	fmt.Println("integration_environment::WithClient -- Creating a new client")
	client, err := createPostgresClient(
		intEnv.Credentials.getConnectionString(),
	)
	if err != nil {
		fmt.Println("integration_environment::WithClient -- Error: failed to create postgres client")
		intEnv.Error = fmt.Errorf("failed to create postgres client: %v", err)
		return intEnv
	}

	intEnv.Client = client

	return intEnv
}

func (intEnv *IntegrationEnvironment) WithMigration() *IntegrationEnvironment {
	if intEnv.Error != nil || intEnv.Client == nil || intEnv.migrated == true {
		return intEnv
	}

	ctx := context.Background()
	client := intEnv.Client
	if client == nil {
		fmt.Println("integration_environment::WithMigration -- Client not found in environment, creating a new client")
		intEnv = intEnv.WithClient()
		client = intEnv.Client
	}

	if err := client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		intEnv.Error = fmt.Errorf("failed to migrate the database: %v", err)
		return intEnv
	}

	fmt.Println("Successfully migrated the Integration database")

	intEnv.migrated = true
	return intEnv
}

func (intEnv *IntegrationEnvironment) Close() {
	if intEnv == nil {
		fmt.Println("integration_environment::Close -- Integration environment not found")
		return
	}

	fmt.Println("integration_environment::Close -- Closing Integration environment")

	if intEnv.Client != nil {
		fmt.Println("integration_environment::Close -- Closing client")
		intEnv.Client.Close()
	}

	if intEnv.Container != nil {
		intEnv.Container.Terminate(context.Background())
	}
}

func getContainerConnectionDelayMs() int {
	containerConnectionDelayMs := 500

	if err := godotenv.Load("../../.env.Integration"); err != nil {
		return containerConnectionDelayMs
	}

	containerConnectionDelayMsStr, ok := os.LookupEnv("CONTAINER_CONNECTION_DELAY_MS")
	if !ok || containerConnectionDelayMsStr == "" {
		return containerConnectionDelayMs
	}

	containerConnectionDelayMs, err := strconv.Atoi(containerConnectionDelayMsStr)
	if err != nil {
		return containerConnectionDelayMs
	}

	if containerConnectionDelayMs < 0 {
		return 0
	}

	return containerConnectionDelayMs
}
