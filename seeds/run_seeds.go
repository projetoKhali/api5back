package seeds

import (
	"fmt"
	"reflect"
	"runtime"

	"api5back/ent"
	"api5back/src/database"
)

func getFuncName(i interface{}) string {
	return runtime.
		FuncForPC(
			reflect.
				ValueOf(i).
				Pointer(),
		).
		Name()
}

func Execute(
	databasePrefix string,
	seedsFunc func(client *ent.Client) error,
) {
	fmt.Printf(
		"Seeding database %s with %s\n",
		databasePrefix,
		getFuncName(seedsFunc),
	)

	client, err := database.Setup(databasePrefix)
	if err != nil {
		panic(fmt.Errorf("failed to setup database: %v", err))
	}
	defer client.Close()

	if err := seedsFunc(client); err != nil {
		panic(fmt.Errorf("failed to seed database: %v", err))
	}

	fmt.Printf(
		"Successfully seeded database %s with %s\n",
		databasePrefix,
		getFuncName(seedsFunc),
	)
}
