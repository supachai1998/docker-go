package main

import (
	"docker_test/routes"
	"docker_test/src/db"
	"docker_test/src/helpers"
	"docker_test/src/logs"
	"log"
)

func main() {
	if err := helpers.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	if err := db.SetupDB(); err != nil {
		log.Fatal(err)
	}

	go logs.WriteLog()

	routes.Setup()
}
