package main

import (
	"fmt"
	"os"
	"strings"

	"entgo.io/contrib/schemast"
)

func createSchemaFile(ctx *schemast.Context, name string) error {
	if err := schemast.Mutate(ctx,
		&schemast.UpsertSchema{
			Name: name,
		},
	); err != nil {
		return fmt.Errorf("scripts/schema • failed to create schema: %v", err)
	}

	if err := ctx.Print("src/schema"); err != nil {
		return fmt.Errorf("scripts/schema • failed to print schema: %v", err)
	}

	return nil
}

func parseName(name string) string {
	// make the first letter uppercase if it is not
	if name[0] >= 'a' && name[0] <= 'z' {

		fmt.Printf("scripts/schema • Warning: schema name `%s` should start with an uppercase letter\n", name)
		name = string(name[0]-32) + name[1:]
	}

	return name
}

func askForOverwrite() (bool, error) {
	var confirm string

	if _, err := fmt.Scan(&confirm); err != nil {
		return false, fmt.Errorf("scripts/schema • failed to read confirmation: %v", err)
	}

	if strings.ToLower(confirm) != strings.ToLower("y") {
		fmt.Println("scripts/schema • Aborted.")
		return false, nil
	}

	return true, nil
}

func CreateSchema(ctx *schemast.Context, name string) error {
	fmt.Printf("scripts/schema • Creating schema `%s`\n", name)

	if err := createSchemaFile(ctx, name); err != nil {
		return fmt.Errorf("scripts/schema • failed to create schema: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		panic("scripts/schema • missing schema name argument")
	}

	ctx, err := schemast.Load("src/schema")
	if err != nil {
		panic(fmt.Errorf("scripts/schema • failed to load context: %v", err))
	}

	name := parseName(os.Args[1])

	// check if the schema already exists, if so, ask for confirmation
	if _, err := os.Stat(fmt.Sprintf("src/schema/%s.go", strings.ToLower(name))); err == nil {
		fmt.Printf("scripts/schema • Schema `%s` already exists, do you want to overwrite it? (y/N): ", name)

		confirm, err := askForOverwrite()
		if err != nil {
			panic(fmt.Errorf("scripts/schema • failed to ask for overwrite: %v", err))
		}

		if !confirm {
			return
		}
	}

	if err := CreateSchema(ctx, name); err != nil {
		panic(fmt.Errorf("scripts/schema • failed to create schema: %v", err))
	}

	fmt.Println("scripts/schema • Successfully created schema.")
}
