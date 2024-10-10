package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsPublic  bool      `json:"is_public"`
}

func ShowCreateNotePage(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}

func CreateNote(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		isPublic := c.PostForm("is_public") == "on"

		_, err := db.Exec(`INSERT INTO notes (title, content, is_public) VALUES (?, ?, ?)`, title, content, isPublic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

func ListNotes(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, content, is_public, created_at FROM notes")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notes"})
			return
		}
		defer rows.Close()

		var notes []Note
		for rows.Next() {
			var note Note
			if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.IsPublic, &note.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while scanning notes"})
				return
			}
			notes = append(notes, note)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Notes": notes,
		})
	}
}

func DeleteNote(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec(`DELETE FROM notes WHERE id = ?`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

func PreviewSharedNote(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var note Note
		query := "SELECT id, title, content, created_at FROM notes WHERE id = ? AND is_public = 1"
		err := db.QueryRow(query, id).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}

		c.HTML(http.StatusOK, "preview.html", note)
	}
}
