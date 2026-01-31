package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var product = []Product{
	{ID: 1, Name: "Copper Ore", Price: 500, Stock: 100},
	{ID: 2, Name: "Silver Ore", Price: 1000, Stock: 70},
}

var category = []Category{
	{ID: 1, Name: "Material", Description: "An item to upgrade tools"},
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func getProductDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	for i := range product {
		if product[i].ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product[i])
			return
		}
	}
}

func addProduct(w http.ResponseWriter, r *http.Request) {
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

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	for i, p := range product {
		if p.ID == id {
			product = append(product[:i], product[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Delete success",
			})

			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	for i := range product {
		if product[i].ID == id {
			var updatedProduct Product
			updatedProduct.ID = product[i].ID

			err := json.NewDecoder(r.Body).Decode(&updatedProduct)
			if err != nil {
				http.Error(w, "ID not found", http.StatusNotFound)
				return
			}

			product[i] = updatedProduct
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product[i])
		}
	}

}

func getCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func getCategoryDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	for i := range category {
		if category[i].ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category[i])
			return
		}
	}
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newCategory.ID = len(category) + 1
	category = append(category, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	for i, c := range category {
		if c.ID == id {
			category = append(category[:i], category[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Delete success",
			})

			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	for i := range category {
		if category[i].ID == id {
			var updatedCategory Category
			updatedCategory.ID = category[i].ID

			err := json.NewDecoder(r.Body).Decode(&updatedCategory)
			if err != nil {
				http.Error(w, "ID not found", http.StatusNotFound)
				return
			}

			category[i] = updatedCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category[i])
		}
	}

}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// GET Product
	// POST Product
	http.HandleFunc("/api/product", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProduct(w, r)
		} else if r.Method == "POST" {
			addProduct(w, r)
		}
	})

	// GET Product Detail
	// DELETE Product
	// PUT Product
	http.HandleFunc("/api/product/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProductDetail(w, r)
		} else if r.Method == "DELETE" {
			deleteProduct(w, r)
		} else if r.Method == "PUT" {
			updateProduct(w, r)
		}
	})

	// GET Category
	// POST Category
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategory(w, r)
		} else if r.Method == "POST" {
			addCategory(w, r)
		}
	})

	// GET Category Detail
	// DELETE Category
	// PUT Category
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryDetail(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
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

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running on ", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Failed to run server")
	}
}
