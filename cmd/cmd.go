package main

import (
	"github.com/vedsatt/calc_online/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
