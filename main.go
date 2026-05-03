package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

// Структуры данных
type Issue struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectID   int    `json:"project_id"`
	AuthorID    int    `json:"author_id"`
	AssigneeID  int    `json:"assignee_id"`
	StatusID    int    `json:"status_id"`
	SeverityID  int    `json:"severity_id"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "./bugtracker.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация БД (все 8 таблиц по списку)
	initDB()

	// Настройка маршрутов
	http.HandleFunc("/api/issues", handleIssues)        // GET, POST
	http.HandleFunc("/api/issues/assign", handleAssign) // PATCH
	http.HandleFunc("/api/stats", handleStats)          // GET (для графиков)

	fmt.Println("Профессиональный Бэкенд запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() {
	queries := []string{
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, role TEXT);",
		"CREATE TABLE IF NOT EXISTS projects (id INTEGER PRIMARY KEY, name TEXT, description TEXT);",
		"CREATE TABLE IF NOT EXISTS statuses (id INTEGER PRIMARY KEY, name TEXT);",
		"CREATE TABLE IF NOT EXISTS severities (id INTEGER PRIMARY KEY, name TEXT);",
		`CREATE TABLE IF NOT EXISTS issues (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			project_id INTEGER REFERENCES projects(id),
			author_id INTEGER REFERENCES users(id),
			assignee_id INTEGER REFERENCES users(id),
			status_id INTEGER REFERENCES statuses(id),
			severity_id INTEGER REFERENCES severities(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		"CREATE TABLE IF NOT EXISTS issue_comments (id INTEGER PRIMARY KEY, issue_id INTEGER, text TEXT);",
		"CREATE TABLE IF NOT EXISTS attachments (id INTEGER PRIMARY KEY, issue_id INTEGER, url TEXT);",
		"CREATE TABLE IF NOT EXISTS project_assignments (id INTEGER PRIMARY KEY, project_id INTEGER, user_id INTEGER);",
		// Наполнение справочников
		"INSERT OR IGNORE INTO statuses (id, name) VALUES (1, 'Open'), (2, 'In Progress'), (3, 'Resolved');",
		"INSERT OR IGNORE INTO severities (id, name) VALUES (1, 'Low'), (2, 'Medium'), (3, 'High'), (4, 'Critical');",
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Printf("Ошибка БД: %v", err)
		}
	}
}

// Обработчики API
func handleIssues(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == http.MethodGet {
		rows, _ := db.Query("SELECT id, title, description, status_id FROM issues")
		defer rows.Close()
		var res []Issue
		for rows.Next() {
			var i Issue
			rows.Scan(&i.ID, &i.Title, &i.Description, &i.StatusID)
			res = append(res, i)
		}
		json.NewEncoder(w).Encode(res)
	} else if r.Method == http.MethodPost {
		var i Issue
		json.NewDecoder(r.Body).Decode(&i)
		db.Exec("INSERT INTO issues (title, description, status_id) VALUES (?, ?, ?)", i.Title, i.Description, 1)
		w.WriteHeader(http.StatusCreated)
	}
}

func handleAssign(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == http.MethodPatch {
		id := r.URL.Query().Get("id")
		user := r.URL.Query().Get("user_id")
		db.Exec("UPDATE issues SET assignee_id = ? WHERE id = ?", user, id)
		fmt.Fprint(w, `{"status":"ok"}`)
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	// Демо-данные для графика менеджера
	stats := map[string]int{"Resolved": 12, "Open": 5}
	json.NewEncoder(w).Encode(stats)
}

func setupCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Content-Type", "application/json")
}
