package main

import (
	"tinyfetch/internal/config"
	"tinyfetch/internal/printer"
)

func main(){
	cfg := config.LoadConfig()

	printer.Print(cfg)
}