package main

import (
        "database/sql"
        "fmt"
        "log"
        "net/http"
        "todo_app/handlers"
        "todo_app/services"

        "github.com/gorilla/mux"
        "github.com/rs/cors"
        _ "github.com/go-sql-driver/mysql"
)

func main() {
        // Initialize database connection
        db, err := sql.Open("mysql", "shubham:root@tcp(mysql:3306)/todo_app")
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
                log.Println("Warning: Error creating todos table (it might already exist):", err)
        }

        // Initialize TodoService
        todoService := &services.TodoServiceImpl{DB: db}

        // Initialize TodoHandler with TodoService
        todoHandler := handlers.NewTodoHandler(todoService)

        // Define routes using Gorilla Mux
        router := mux.NewRouter()

        router.HandleFunc("/todos", todoHandler.GetAllTodos).Methods("GET")
        router.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
        router.HandleFunc("/todos/{id}", todoHandler.GetTodoByID).Methods("GET")
        router.HandleFunc("/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
        router.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

        // CORS middleware
        corsHandler := cors.Default().Handler(router)

        // Start the server
        port := ":8082"
        fmt.Println("Server started on port", port)
        log.Fatal(http.ListenAndServe("0.0.0.0"+port, corsHandler))
}                              
