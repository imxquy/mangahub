package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"mangahub/pkg/database"
)

func main() {
	dbPath := getenv("MANGAHUB_DB", "./mangahub.db")
	jwtSecret := getenv("MANGAHUB_JWT_SECRET", "dev-secret")

	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	r := gin.Default()

	_ = jwtSecret
	_ = db

	log.Println("HTTP API listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
