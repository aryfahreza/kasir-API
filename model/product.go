package model

type Product struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	Stock         int    `json:"stock"`
	ID_Category   int    `json:"id_category"`
	Category_Name string `json:"category_name"`
}
