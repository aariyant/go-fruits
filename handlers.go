package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	baseTemplate := "src/templates/base.html"
	tmpl = fmt.Sprintf("src/templates/%s", tmpl)
	t, err := template.ParseFiles(baseTemplate, tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	// Render index template
	renderTemplate(w, "index.html", nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	// Query the database for all fruits
	rows, err := DB.Query("SELECT * FROM fruits")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var fruits []Fruit
	for rows.Next() {
		var fruit Fruit
		if err := rows.Scan(&fruit.ID, &fruit.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fruits = append(fruits, fruit)
	}

	// Render list template with the list of fruits
	// Render the list.html template with the list of fruits
	renderTemplate(w, "list.html", fruits)
}

func fruitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the fruit ID from the URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Fruit ID not provided", http.StatusBadRequest)
		return
	}

	// Query the database to retrieve the fruit details by ID
	var fruit Fruit
        query := fmt.Sprintf("SELECT id, name FROM fruits WHERE id = '%s'", id)
	err := DB.QueryRow(query).Scan(&fruit.ID, &fruit.Name)
	if err != nil {
		http.Error(w, "Fruit not found", http.StatusNotFound)
		return
	}

	// Render the single.html template with the fruit details
	renderTemplate(w, "fruit.html", fruit)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Render the add.html template for GET requests
		renderTemplate(w, "add.html", nil)

	case http.MethodPost:

		name := r.FormValue("name")

		// Insert new fruit into the database
		query := fmt.Sprintf("INSERT INTO fruits (name) VALUES ('%s')", name)
		_, err := DB.Query(query)
		if err != nil {
			http.Error(w, "Failed to insert into the database", http.StatusInternalServerError)
			log.Printf("Failed to insert into the database: %v", err)
			log.Printf(query)
			return
		}

		http.Redirect(w, r, "/list", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func editHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Render the edit.html template for GET requests
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Fruit ID not provided", http.StatusBadRequest)
			return
		}

		// Query the database to retrieve the fruit details by ID
		var fruit Fruit
                query := fmt.Sprintf("SELECT id, name FROM fruits WHERE id = '%s'", id)
		err := DB.QueryRow(query).Scan(&fruit.ID, &fruit.Name)
		if err != nil {
			http.Error(w, "Fruit not found", http.StatusNotFound)
			return
		}

		renderTemplate(w, "edit.html", fruit)
	case http.MethodPost:
		// Handle form submission and update the fruit for POST requests
		id := r.FormValue("id")
		name := r.FormValue("name")

		query := fmt.Sprintf("UPDATE fruits SET name = '%s' WHERE id = '%s'", name, id)
		_, err := DB.Exec(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/list", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Render the delete.html template for GET requests
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Fruit ID not provided", http.StatusBadRequest)
			return
		}

		// Query the database to retrieve the fruit details by ID
		var fruit Fruit
                query := fmt.Sprintf("SELECT id, name FROM fruits WHERE id = '%s'", id)
		err := DB.QueryRow(query).Scan(&fruit.ID, &fruit.Name)
		if err != nil {
			http.Error(w, "Fruit not found", http.StatusNotFound)
			return
		}

		renderTemplate(w, "delete.html", fruit)
	case http.MethodPost:
		// Handle form submission and delete the fruit for POST requests
		id := r.FormValue("id")

		query := fmt.Sprintf("DELETE FROM fruits WHERE id = '%s'", id)
		_, err := DB.Exec(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/list", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
