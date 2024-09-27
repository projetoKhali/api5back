package main

import (
	"fmt"
	"os"
	"strings"

	"api5back/ent"
	"api5back/seeds"
)

type SeedsPreset struct {
	Abbreviation string
	Name         string
	SeedsFunc    func(client *ent.Client) error
}

var SeedsPresets = []SeedsPreset{
	{
		Abbreviation: "dw",
		Name:         "DataWarehouse",
		SeedsFunc:    seeds.DataWarehouse,
	},
	{
		Abbreviation: "phpc",
		Name:         "ProceduralHiringProcessCandidates",
		SeedsFunc:    seeds.DwProceduralHiringProcessCandidates,
	},
}

func buildAvailableSeedsMessage(sb *strings.Builder) {
	sb.WriteString("Available seeds:\n\n")
	sb.WriteString("Abbreviation\tName\n")
	for _, seed := range SeedsPresets {
		sb.WriteString(fmt.Sprintf(
			"• %s\t\t• %s\n",
			seed.Abbreviation,
			seed.Name,
		))
	}
	sb.WriteString("\n")
}

func askForConfirmation() (bool, error) {
	var confirm string

	if _, err := fmt.Scan(&confirm); err != nil {
		return false, fmt.Errorf("scripts/seeds • failed to read confirmation: %v", err)
	}

	if strings.ToLower(confirm) != strings.ToLower("y") {
		fmt.Println("scripts/seeds • Aborted.")
		return false, nil
	}

	return true, nil
}

// Ponto de entrada manual para rodar a seed
func main() {
	// parse args here (could be name or abbreviation)
	// and run the appropriate SeedsFunc from the map
	if len(os.Args) < 2 {
		sb := strings.Builder{}
		sb.WriteString("scripts/seeds • missing seed name argument\n")
		buildAvailableSeedsMessage(&sb)
		sb.WriteString("You can run a seed by providing its name or abbreviation as an argument.\n")
		sb.WriteString("Example:\n")
		sb.WriteString("  go run scripts/seeds.go DataWarehouse\n")
		sb.WriteString("  go run scripts/seeds.go dw\n")
		panic(sb.String())
	}

	seedName := os.Args[1]
	var targetSeedsPreset *SeedsPreset
	for _, seedsPreset := range SeedsPresets {
		if strings.ToLower(seedsPreset.Name) == strings.ToLower(seedName) ||
			strings.ToLower(seedsPreset.Abbreviation) == strings.ToLower(seedName) {
			targetSeedsPreset = &seedsPreset
			break
		}
	}

	if targetSeedsPreset == nil {
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("scripts/seeds • seed `%s` not found\n", seedName))
		buildAvailableSeedsMessage(&sb)
		panic(sb.String())
	}

	// ask for confirmation
	fmt.Printf(
		"scripts/seeds • Running seed `%s`, confirm? (y/N): ",
		targetSeedsPreset.Name,
	)

	confirm, err := askForConfirmation()
	if err != nil {
		panic(fmt.Errorf("scripts/seeds • failed to ask for confirmation: %v", err))
	}

	if !confirm {
		return
	}

	seeds.Execute("DW", targetSeedsPreset.SeedsFunc)
}
