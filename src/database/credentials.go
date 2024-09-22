package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseCredentials struct {
	User string
	Pass string
	Host string
	Port *int
	Name *string
}

func internalNewDatabaseCredentials(prefix string, envFilePath string) (*DatabaseCredentials, error) {
	fields := []string{
		fmt.Sprintf("%s_USER", prefix),
		fmt.Sprintf("%s_PASS", prefix),
		fmt.Sprintf("%s_HOST", prefix),
		fmt.Sprintf("%s_PORT", prefix),
		fmt.Sprintf("%s_NAME", prefix),
	}

	err := godotenv.Load(envFilePath)
	if err != nil {
		return nil, fmt.Errorf(
			"error loading `%s` file: %w",
			envFilePath,
			err,
		)
	}

	for i, field := range fields {
		value, ok := os.LookupEnv(field)
		if !ok {
			return nil, fmt.Errorf(
				"missing environment variable `%s` in `%s` file",
				field,
				envFilePath,
			)
		}

		fields[i] = value
	}

	port, err := strconv.Atoi(fields[3])
	if err != nil {
		return nil, fmt.Errorf(
			"invalid port number: `%s` in `%s`: %w",
			fields[3],
			envFilePath,
			err,
		)
	}

	return &DatabaseCredentials{
		User: fields[0],
		Pass: fields[1],
		Host: fields[2],
		Port: &port,
		Name: &fields[4],
	}, nil
}

func newDatabaseCredentials(prefix string) (*DatabaseCredentials, error) {
	return internalNewDatabaseCredentials(
		prefix,
		".env",
	)
}

func newTestingDatabaseCredentials() (*DatabaseCredentials, error) {
	return internalNewDatabaseCredentials(
		"TI",
		"../../.env.integration",
	)
}

func (dc *DatabaseCredentials) getConnectionString() string {
	connectionString := fmt.Sprintf(
		"user=%s password=%s host=%s",
		dc.User,
		dc.Pass,
		dc.Host,
	)

	if dc.Port != nil && *dc.Port != 0 {
		connectionString = fmt.Sprintf(
			"%s port=%d",
			connectionString,
			*dc.Port,
		)
	}

	if dc.Name != nil && *dc.Name != "" {
		connectionString = fmt.Sprintf(
			"%s dbname=%s",
			connectionString,
			*dc.Name,
		)
	}

	return fmt.Sprintf(
		"%s sslmode=disable",
		connectionString,
	)
}
