package main

import (
	"github.com/amavrin/go-musthave-devops/internal/client"
	"github.com/amavrin/go-musthave-devops/internal/metrics"
)

const (
	serverURL = "http://localhost:8080"
)

func initDB() *metrics.DB {
	db := metrics.NewDB()
	db.RunUpdates()
	return db
}

func initClient() *client.Client {
	return client.NewClient(serverURL)
}

func run() {
	db := initDB()
	client := initClient()
	client.SendLoop(db)
}

func main() {
	run()
}
