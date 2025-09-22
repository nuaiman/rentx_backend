package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	db, err := sql.Open("sqlite", "api.db")
	if err != nil {
		log.Fatal("Database could not connect: ", err)
	}
	DB = db

	if _, err = DB.Exec("PRAGMA busy_timeout = 5000;"); err != nil {
		log.Fatal("Could not set busy_timeout: ", err)
	}

	if _, err = DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatal("Could not enable foreign keys: ", err)
	}

	if _, err = DB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		log.Fatal("Could not enable WAL: ", err)
	}

	if _, err = DB.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		log.Fatal("Could not set synchronous mode: ", err)
	}

	if err := createTables(); err != nil {
		log.Fatal("Database table creation failed: ", err)
	}

	createDefaultSuperAdmin()

	fmt.Println("✅ Database initialized, tables created and superadmin ensured!")
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			fmt.Println("Error closing database:", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}
}

func createDefaultSuperAdmin() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'superadmin'").Scan(&count)
	if err != nil {
		log.Fatalf("❌ Failed to check superadmin count: %v", err)
	}

	if count == 0 {
		name := os.Getenv("SUPERADMIN_NAME")
		email := os.Getenv("SUPERADMIN_EMAIL")
		phone := os.Getenv("SUPERADMIN_PHONE")
		password := os.Getenv("SUPERADMIN_PASSWORD")

		// fail hard if env variables missing
		if name == "" || email == "" || phone == "" || password == "" {
			log.Fatal("❌ SUPERADMIN environment variables not set. Exiting...")
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("❌ Failed to hash password: %v", err)
		}

		_, err = DB.Exec(`
            INSERT INTO users (name, email, phone, password, image, role)
            VALUES (?, ?, ?, ?, ?, 'superadmin')`,
			name, email, phone, string(hashed), "",
		)
		if err != nil {
			log.Fatalf("❌ Failed to insert default superadmin: %v", err)
		}

		fmt.Println("✅ Default superadmin user created from environment variables.")
	}
}

func createTables() error {
	tables := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			phone TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			image TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user', -- 'superadmin' | 'admin' | 'user'
			dateTime DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Refresh-Token table
		`CREATE TABLE IF NOT EXISTS refreshTokens (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		userId INTEGER NOT NULL,
    		token TEXT NOT NULL UNIQUE,
    		expiresAt DATETIME NOT NULL,
			dateTime DATETIME DEFAULT CURRENT_TIMESTAMP,
    		FOREIGN KEY(userId) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// Categories table
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			dateTime DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Posts table
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL,
			categoryId INTEGER NOT NULL,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			description TEXT NOT NULL,
			dailyPrice REAL NOT NULL,
			weeklyPrice REAL NOT NULL,
			monthlyPrice REAL NOT NULL,
			dateTime DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userId) REFERENCES users (id) ON DELETE CASCADE,
			FOREIGN KEY (categoryId) REFERENCES categories (id) ON DELETE CASCADE
		)`,

		// Post images table (normalized)
		`CREATE TABLE IF NOT EXISTS post_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			postId INTEGER NOT NULL,
			imageUrl TEXT NOT NULL,
			position INTEGER DEFAULT 0,
			FOREIGN KEY (postId) REFERENCES posts (id) ON DELETE CASCADE
		)`,

		// Orders table
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL,
			postId INTEGER NOT NULL,
			dateTime DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userId) REFERENCES users (id) ON DELETE CASCADE,
			FOREIGN KEY (postId) REFERENCES posts (id) ON DELETE CASCADE
		)`,

		// Reviews table
		`CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL,
			postId INTEGER NOT NULL,
			review TEXT NOT NULL,
			dateTime DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userId) REFERENCES users (id) ON DELETE CASCADE,
			FOREIGN KEY (postId) REFERENCES posts (id) ON DELETE CASCADE
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_posts_categoryId ON posts(categoryId)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_userId ON orders(userId)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_postId ON orders(postId)`,
		`CREATE INDEX IF NOT EXISTS idx_reviews_userId ON reviews(userId)`,
		`CREATE INDEX IF NOT EXISTS idx_reviews_postId ON reviews(postId)`,
		`CREATE INDEX IF NOT EXISTS idx_post_images_postId ON post_images(postId)`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return fmt.Errorf("error creating table/index: %w", err)
		}
	}

	return nil
}
