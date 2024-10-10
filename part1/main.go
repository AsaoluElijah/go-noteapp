package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a table
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	// Insert a user
	_, err = db.Exec("INSERT INTO users(name) VALUES(?)", "John Doe")
	if err != nil {
		log.Fatal(err)
	}

	// Query users
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("User: %d, %s", id, name)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// UPDATE
	result, err := db.Exec("UPDATE users SET name = ? WHERE id = ?", "Jane Doe", 3)
	if err != nil {
		log.Fatal(err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Updated %d row(s)\n", affectedRows)

	// DELETE
	result, err = db.Exec("DELETE FROM users WHERE id = ?", 4)
	if err != nil {
		log.Fatal(err)
	}

	affectedRows, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted %d row(s)\n", affectedRows)
}
