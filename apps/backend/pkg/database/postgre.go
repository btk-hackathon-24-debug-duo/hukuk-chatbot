package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func SetupDb() (*sql.DB, error) {
	config := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432, // Default to 5432 if DB_PORT is not set
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SLL_MODE"),
	}

	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		} else {
			fmt.Printf("Warning: Invalid DB_PORT value '%s', using default 5432\n", portStr)
		}
	}

	fmt.Printf("Attempting to connect to database with config: %+v\n", config)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	err = SetupTables(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("failed to setup tables: %w", err)
	}

	fmt.Println("Successfully connected to the database")

	return sqlDB, nil

}

func SetupTables(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS users ( 
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil

}
