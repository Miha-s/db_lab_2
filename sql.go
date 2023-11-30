package main

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
)

type UserInsertRequest struct {
	Count int `json:"count"`
}

func registerPostgresHandlers(router *http.ServeMux, pool *pgxpool.Pool) {
	router.HandleFunc("/test/postgres/emploees/insert", 
	func(w http.ResponseWriter, r *http.Request) {
		insertTestEmploeesPostgresHandler(w, r, pool)
	})
	router.HandleFunc("/test/postgres/emploees/deleteAll", 
	func(w http.ResponseWriter, r *http.Request) {
		deleteAllEmploeeysDataPostgresHandler(w, r, pool)
	})
	router.HandleFunc("/test/postgres/emploees/update", 
	func(w http.ResponseWriter, r *http.Request) {
		updateAllEmploeesPositionPostgresHandler(w, r, pool)
	})
	router.HandleFunc("/test/postgres/emploees/add_skills", 
	func(w http.ResponseWriter, r *http.Request) {
		addSkillsForEmploeesPostgresHandler(w, r, pool)
	})
}

func addSkillsForEmploeesPostgresHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()

	ctx := context.Background()

	rows, err := pool.Query(context.Background(),
		"SELECT id FROM employees",
	)
	defer rows.Close()

	for rows.Next() {
		var employeeID int
		if err := rows.Scan(&employeeID); err != nil {
			log.Printf("Error: %v\n", err)
		}

		_, err = pool.Exec(ctx,
			"INSERT INTO employee_skills (employee_id, skill) VALUES ($1, $2)",
			employeeID, "new skill",
		)
	}

	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}


	duration := time.Since(start).Seconds()
	sendResponse(w, "PostgreSQL", duration)
}

func updateAllEmploeesPositionPostgresHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()

	_, err := pool.Exec(context.Background(),
		"UPDATE employees SET current_position = $1",
		"Update position",
	)

	if err != nil {
		http.Error(w, "Error updating user passwords", http.StatusInternalServerError)
		log.Printf("Error updating user passwords: %v\n", err)
		return
	}

	duration := time.Since(start).Seconds()
	sendResponse(w, "PostgreSQL", duration)
}

func deleteAllEmploeeysDataPostgresHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	tables := []string{"employee_projects", "employee_skills", "employees"}
	for _, table := range tables {
		if _, err := tx.Exec(context.Background(), "DELETE FROM "+table); err != nil {
			http.Error(w, "Error deleting data from table "+table, http.StatusInternalServerError)
			tx.Rollback(context.Background())
			return
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	log.Println("All user data deleted successfully")
	duration := time.Since(start).Seconds()
	sendResponse(w, "PostgreSQL", duration)
}

func insertTestEmploeesPostgresHandler(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UserInsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Count <= 0 {
		http.Error(w, "Count must be a positive integer", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	start := time.Now()
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	for i := 0; i < req.Count; i++ {
		var employeeID int
		err = tx.QueryRow(ctx,
			"INSERT INTO employees (name, surname, email, telephone, employment_date, current_position, current_salary) VALUES ('John', 'Doe', 'john.doe@example.com', '1234567890', NOW(), 'Software Engineer', 75000.0) RETURNING id",
		).Scan(&employeeID)
	
		if err != nil {
			log.Fatal("Error inserting employee:", err)
		}

		_, err = tx.Exec(ctx,
		"INSERT INTO employee_skills (employee_id, skill) VALUES ($1, $2)",
		employeeID, "some skill",
		)

		if err != nil {
			log.Fatal("Error inserting skill:", err)
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO employee_projects (employee_id, project_id) VALUES ($1, $2)",
			employeeID, 1,
		)
		if err != nil {
			log.Fatal("Error inserting data into employee_projects:", err)
		}
	
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	sendResponse(w, "PostgreSQL", duration)
}