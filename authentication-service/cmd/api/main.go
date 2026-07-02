package main

import (
	"authentication/cmd/api/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "8080"

type Application struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service...")
	db := connectToDB()
	if db == nil {
		log.Fatal("Failed to connect to database")
	}
	app := Application{
		DB:     db,
		Models: data.New(db),
	}

	server := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	if dsn == "" {
		log.Fatal("DSN environment variable is not set")
	}

	log.Println("Connecting to database with DSN:", dsn)

	for i := 1; i <= 10; i++ {
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			log.Printf("Failed to create DB handle: %v", err)
		} else {
			err = db.Ping()
			if err == nil {
				log.Println("Connected to database")
				return db
			}

			log.Printf("Database ping failed: %v", err)
			db.Close()
		}

		log.Printf("Waiting for database to connect... attempt %d/10", i)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Failed to connect to database after 10 attempts")
	return nil
}
