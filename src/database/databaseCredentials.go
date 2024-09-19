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

func NewDatabaseCredentials() (*DatabaseCredentials, error) {
	fields := []string{
		"DB_USER",
		"DB_PASS",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	}

	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading `.env` file: %w", err)
	}

	for i, field := range fields {
		value, ok := os.LookupEnv(field)
		if !ok {
			return nil, fmt.Errorf("missing environment variable `%s` in `.env` file", field)
		}

		fields[i] = value
	}

	port, err := strconv.Atoi(fields[3])
	if err != nil {
		return nil, fmt.Errorf("invalid port number: `%s` in `.env` file: %w", fields[3], err)
	}

	return &DatabaseCredentials{
		User: fields[0],
		Pass: fields[1],
		Host: fields[2],
		Port: &port,
		Name: &fields[4],
	}, nil
}

func (dc *DatabaseCredentials) GetConnectionString() string {
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
