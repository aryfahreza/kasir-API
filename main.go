package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type Categories struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var product = []Product{
	{ID: 1, Name: "Copper Ore", Price: 500, Stock: 100},
	{ID: 2, Name: "Silver Ore", Price: 1000, Stock: 70},
}

var categories = []Categories{
	{ID: 1, Name: "Material", Description: "An item to upgrade tools"},
}

func main() {
	// localhost:8080/api/product
	http.HandleFunc("/api/product", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
		} else if r.Method == "POST" {
			var newProduct Product
			err := json.NewDecoder(r.Body).Decode(&newProduct)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			newProduct.ID = len(product) + 1
			product = append(product, newProduct)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newProduct)
		}
	})

	// GET Categories
	// POST Categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			var newCategory Categories

			err := json.NewDecoder(r.Body).Decode(&newCategory)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			newCategory.ID = len(categories) + 1
			categories = append(categories, newCategory)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(categories)
		}
	})

	// GET Category Detail
	// PUT Category
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "ID not found", http.StatusNotFound)
			return
		}

		for i := range categories {
			if categories[i].ID == id {
				if r.Method == "GET" {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(categories[i])
					return
				} else if r.Method == "PUT" {
					var updatedCategory Categories
					updatedCategory.ID = categories[i].ID

					err := json.NewDecoder(r.Body).Decode(&updatedCategory)
					if err != nil {
						http.Error(w, "ID not found", http.StatusNotFound)
						return
					}

					categories[i] = updatedCategory
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(categories[i])
				}
			}
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running on localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to run server")
	}
}
