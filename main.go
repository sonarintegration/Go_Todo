package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"todo_app/handlers"
	"todo_app/services"

	"github.com/rs/cors"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize database connection
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/todo_app")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Create todos table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		completed TINYINT(1) NOT NULL DEFAULT 0
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error creating todos table:", err)
	}

	// Initialize TodoService
	todoService := &services.TodoServiceImpl{DB: db}

	// Initialize TodoHandler with TodoService
	todoHandler := handlers.NewTodoHandler(todoService)

	// Define routes
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoHandler.GetAllTodos(w, r)
		case http.MethodPost:
			todoHandler.CreateTodo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			todoHandler.GetTodoByID(w, r)
		} else if r.Method == http.MethodPut {
			todoHandler.UpdateTodo(w, r)
		} else if r.Method == http.MethodDelete {
			todoHandler.DeleteTodo(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// CORS middleware
	corsHandler := cors.Default().Handler(mux)

	// Start the server
	port := ":8080"
	fmt.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe("172.27.59.220"+port, corsHandler))
}
