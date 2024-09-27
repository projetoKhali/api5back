package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Credentials struct {
	User string
	Pass string
	Host string
	Port *int
	Name *string
}

func internalNewCredentials(prefix string, envFilePath string) (*Credentials, error) {
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

	return &Credentials{
		User: fields[0],
		Pass: fields[1],
		Host: fields[2],
		Port: &port,
		Name: &fields[4],
	}, nil
}

func newCredentials(prefix string) (*Credentials, error) {
	return internalNewCredentials(
		prefix,
		".env",
	)
}

func newTestingCredentials() (*Credentials, error) {
	return internalNewCredentials(
		"TI",
		"../../.env.integration",
	)
}

func (dc *Credentials) getConnectionString() string {
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
