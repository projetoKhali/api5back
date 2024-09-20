package main

import (
	"fmt"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

// defines the settings for the ent codegen
func entConfig() *gen.Config {
	return &gen.Config{
		Schema:  "database/schema",
		Target:  "ent",
		Package: "api5back/ent",
	}
}

// public function to generate ent files
func GenerateEntFiles() error {
	if err := entc.Generate("./src/schema", entConfig()); err != nil {
		return fmt.Errorf("scripts/generate • error running ent codegen: %w", err)
	}

	return nil
}

// manual entry point for codegen on command
func main() {
	fmt.Println("scripts/generate • Generating ent files...")

	if err := GenerateEntFiles(); err != nil {
		panic(fmt.Errorf("scripts/generate • failed to generate ent files: %w", err))
	}

	fmt.Println("scripts/generate • Successfully genereted ent files.")
}
