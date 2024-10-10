package main

import (
	"database/sql"
	"log"

	"go-sqlite/handlers" // ⚠️ Change this line

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		is_public INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/", handlers.ListNotes(db))
	router.GET("/create", handlers.ShowCreateNotePage)
	router.POST("/create", handlers.CreateNote(db))
	router.POST("/notes/:id/delete", handlers.DeleteNote(db))
	router.GET("/share/:id", handlers.PreviewSharedNote(db))

	port := ":8080"
	log.Printf("Server is running on http://localhost%s", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
