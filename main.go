package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"sqlboiler-project/models"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=sqlboiler_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	book := models.Book{
		Title:  "Sample Book",
		Author: "John Doe",
	}

	err = book.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Book inserted successfully")
}
