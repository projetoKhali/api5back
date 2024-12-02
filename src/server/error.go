package server

import (
	"os"
)

func DisplayError(err error) string {
	if os.Getenv("LOCALHOST") == "TRUE" {
		return err.Error()
	} else {
		return "Erro"
	}
}
