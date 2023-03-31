package main

import (
	_ "github.com/lib/pq"

	"tg_service/internal/app"
)

const serviceName = "tg_service"

func main() {
	a := app.New(serviceName)
	a.Run()
}
