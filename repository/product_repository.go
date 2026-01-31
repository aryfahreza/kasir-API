package repository

import (
	"database/sql"
	"errors"
	"kasir-api/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]model.Product, error) {
	query := "SELECT id, name, price, stock FROM product"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	product := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}

		product = append(product, p)
	}

	return product, nil
}

func (repo *ProductRepository) Create(product *model.Product) error {
	query := "INSERT INTO product (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) GetByID(id int) (*model.Product, error) {
	query := "SELECT id, name, price, stock FROM product WHERE id = $1"

	var p model.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)

	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}
