package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "demo"
)

func main() {
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping", err)
	}

	fmt.Println("successfully connected")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("failed to begin transaction", err)
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			log.Fatal("transaction rolled back due to error:", err)
		}
	}()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			category VARCHAR(255) NOT NULL,
			price INT
		)
	`
	_, err = tx.Exec(createTableSQL)
	if err != nil {
		log.Fatal("failed to create table", err)
	}

	productName := "New Product"
	categoryName := "New Category"
	insertSQL := `
		INSERT INTO products (name, category) VALUES ($1, $2)
	`
	_, err = tx.Exec(insertSQL, productName, categoryName)
	if err != nil {
		log.Fatal("failed to insert product", err)
	}

	newPrice := 100
	updateSQL := `
		UPDATE products SET price = $1 WHERE name = $2
	`
	_, err = tx.Exec(updateSQL, newPrice, productName)
	if err != nil {
		log.Fatal("failed to update product price", err)
	}

	retrieveSQL := `
		SELECT price, category FROM products WHERE name = $1
	`
	var retrievedPrice int
	var retrievedCategory string
	err = tx.QueryRow(retrieveSQL, productName).Scan(&retrievedPrice, &retrievedCategory)
	if err != nil {
		log.Fatal("failed to retrieve product price", err)
	}
	fmt.Println("Retrieved product price and category:", retrievedPrice, retrievedCategory)

	deleteSQL := `
		DELETE FROM products WHERE name = $1
	`
	_, err = tx.Exec(deleteSQL, productName)
	if err != nil {
		log.Fatal("failed to delete product", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("failed to commit", err)
	}

	fmt.Println("Transaction successfully committed")
}
