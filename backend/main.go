package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"log"
)

func main() {
	dbConn := flag.String("db",
		"host=localhost dbname=lamoda user=postgres password=lamoda sslmode=disable",
		"database connection string")
	httpPort := flag.Int("http-port", 4000, "HTTP API port")
	flag.Parse()

	postgres, err := sql.Open("postgres", *dbConn)
	if err != nil {
		log.Fatal(err)
	}

	defer func(postgres *sql.DB) {
		err = postgres.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(postgres)

	db := goqu.New("postgres", postgres)

	s := newHttpServer(db)

	err = s.run(fmt.Sprintf(":%d", *httpPort))
	if err != nil {
		log.Fatal(err)
	}
}
