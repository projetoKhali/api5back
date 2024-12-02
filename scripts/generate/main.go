package main

import (
	"fmt"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

// defines the settings for the ent codegen
func entConfig() *gen.Config {
	return &gen.Config{
		Schema:  "src/schema",
		Target:  "ent",
		Package: "api5back/ent",
		Features: []gen.Feature{
			gen.FeatureModifier,
			gen.FeatureExecQuery,
		},
	}
}

// public function to generate ent files
func GenerateEntFiles() error {
	if err := entc.Generate("./src/schema", entConfig()); err != nil {
		return fmt.Errorf("scripts/generate • error running ent codegen: %w", err)
	}

	return nil
}

func Generate() {
	fmt.Println("scripts/generate • Generating ent files...")

	if err := GenerateEntFiles(); err != nil {
		panic(fmt.Errorf("scripts/generate • failed to generate ent files: %w", err))
	}

	fmt.Println("scripts/generate • Successfully generated ent files.")
}

// manual entry point for codegen on command
func main() {
	Generate()
}
