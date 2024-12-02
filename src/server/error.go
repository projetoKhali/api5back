package server

import (
	"fmt"
	"os"
)

func DisplayError(err error) string {
	if os.Getenv("LOCALHOST") == "TRUE" {
		return err.Error()
	} else {
		fmt.Println(err.Error())
		return "Erro"
	}
}
