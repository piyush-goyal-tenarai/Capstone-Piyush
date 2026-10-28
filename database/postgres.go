package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"order-management/config"
)

func Connect(cfg *config.Config) *sql.DB {
	// Construct a key-value DSN that safely encloses password in quotes
dsn := "host=localhost port=5432 user=postgres password='Piyushgoyal@2005' dbname=capstonedb sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	createTable(db)
	fmt.Println("[DB] Connected to PostgreSQL")
	return db
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL,
		product_id VARCHAR(50) NOT NULL,
		quantity INT NOT NULL CHECK (quantity > 0),
		price DECIMAL(10,2) NOT NULL DEFAULT 0,
		total_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
		status VARCHAR(20) NOT NULL DEFAULT 'CREATED'
			CHECK (status IN ('CREATED', 'CONFIRMED', 'CANCELLED')),
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
	CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
	`
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
