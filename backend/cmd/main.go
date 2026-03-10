package main

import (
	"backend/internal/infrastructure/db"
	"backend/internal/interfaces/http"
)

func main() {
	db.Init()
	defer db.Close()

	r := http.SetupRouter()

	r.Run(":8080")
}
