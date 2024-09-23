package main

import (
	"api5back/seeds"
)

// Ponto de entrada manual para rodar a seed
func main() {
	seeds.Execute("DW", seeds.DataWarehouse)
}
