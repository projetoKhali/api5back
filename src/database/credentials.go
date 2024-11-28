package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Field struct with generics for dynamic typing
type Field[T any] struct {
	Name     string // The name of the parameter
	Value    T      // The value of the parameter
	Required bool   // Whether the parameter is required
}

// Credentials struct using Field types
type Credentials struct {
	User    Field[string]
	Pass    Field[string]
	Host    Field[string]
	Port    Field[int]
	Name    Field[string]
	SSLMode Field[string]
}

// Helper function to load and parse environment variables
func loadField[T any](field *Field[T], parser func(string) (T, error), envFilePath string) error {
	envValue, ok := os.LookupEnv(field.Name)
	if field.Required && !ok {
		return fmt.Errorf("missing required environment variable `%s` in `%s` file", field.Name, envFilePath)
	}

	if !ok || envValue == "" {
		var zeroValue T
		field.Value = zeroValue
		return nil
	}

	parsedValue, err := parser(envValue)
	if err != nil {
		return fmt.Errorf("invalid value for `%s`: `%s`: %w", field.Name, envValue, err)
	}
	field.Value = parsedValue
	return nil
}

func newCredentials(prefix string) (*Credentials, error) {
	return internalNewCredentials(prefix, ".env")
}

func newTestingCredentials() (*Credentials, error) {
	return internalNewCredentials("TI", "../../.env.integration")
}

func internalNewCredentials(prefix string, envFilePath string) (*Credentials, error) {
	// Load the environment file
	if err := godotenv.Load(envFilePath); err != nil {
		return nil, fmt.Errorf("error loading `%s` file: %w", envFilePath, err)
	}

	// Define fields dynamically using the Field struct
	fields := []interface{}{
		&Field[string]{Name: fmt.Sprintf("%s_USER", prefix), Required: true},
		&Field[string]{Name: fmt.Sprintf("%s_PASS", prefix), Required: true},
		&Field[string]{Name: fmt.Sprintf("%s_HOST", prefix), Required: true},
		&Field[int]{Name: fmt.Sprintf("%s_PORT", prefix), Required: true},
		&Field[string]{Name: fmt.Sprintf("%s_NAME", prefix), Required: false},
		&Field[string]{Name: "SSLMODE", Required: false},
	}

	// Define parsers for each field type
	stringParser := func(value string) (string, error) { return value, nil }
	intParser := func(value string) (int, error) { return strconv.Atoi(value) }
	boolParser := func(value string) (bool, error) { return value == "true", nil }

	// Map field types to their respective parsers
	parsers := map[string]interface{}{
		"string": stringParser,
		"int":    intParser,
		"bool":   boolParser,
	}

	// Iterate over fields and populate their values
	for _, field := range fields {
		switch f := field.(type) {
		case *Field[string]:
			if err := loadField(f, parsers["string"].(func(string) (string, error)), envFilePath); err != nil {
				return nil, err
			}
		case *Field[int]:
			if err := loadField(f, parsers["int"].(func(string) (int, error)), envFilePath); err != nil {
				return nil, err
			}
		case *Field[bool]:
			if err := loadField(f, parsers["bool"].(func(string) (bool, error)), envFilePath); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported field type: %T", f)
		}
	}

	// Construct the Credentials struct using the populated fields
	return &Credentials{
		User:    *fields[0].(*Field[string]),
		Pass:    *fields[1].(*Field[string]),
		Host:    *fields[2].(*Field[string]),
		Port:    *fields[3].(*Field[int]),
		Name:    *fields[4].(*Field[string]),
		SSLMode: *fields[5].(*Field[string]),
	}, nil
}

func (dc *Credentials) getConnectionString() string {
	connectionString := fmt.Sprintf(
		"user=%s password=%s host=%s",
		dc.User.Value,
		dc.Pass.Value,
		dc.Host.Value,
	)

	if dc.Port.Value != 0 {
		connectionString = fmt.Sprintf(
			"%s port=%d",
			connectionString,
			dc.Port.Value,
		)
	}

	if dc.Name.Value != "" {
		connectionString = fmt.Sprintf(
			"%s dbname=%s",
			connectionString,
			dc.Name.Value,
		)
	}

	return fmt.Sprintf(
		"%s sslmode=%s",
		connectionString,
		dc.SSLMode.Value,
	)
}
