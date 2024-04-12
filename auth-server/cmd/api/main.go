package main

import (
	"auth-server/data"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	conn, err := connectToDB()
	if err != nil {
		log.Panic(err)
	}
	if conn == nil {
		log.Panic(errors.New("could not connect to db"))
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Printf("Could not connect to DB: %v\n", err)
			counts++
		} else {
			log.Println("Connected to DB")
			return conn, err
		}

		if counts > 10 {
			return nil, errors.New("could not connect to db")
		}

		log.Println("Backing off for 2 secconds")
		time.Sleep(2 * time.Second)
		continue
	}

}
