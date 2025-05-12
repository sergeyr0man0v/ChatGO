package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(config *Config) (*Database, error) {
	var db *sql.DB
	var err error
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
		)

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			fmt.Printf("Failed to open database connection (attempt %d/%d): %v\n", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		err = db.Ping()
		if err != nil {
			fmt.Printf("Failed to ping database (attempt %d/%d): %v\n", i+1, maxRetries, err)
			db.Close()
			time.Sleep(retryDelay)
			continue
		}

		// Connection successful
		return &Database{db: db}, nil
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
