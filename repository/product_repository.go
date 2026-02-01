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
	query := `SELECT p.id, p.name, p.price, p.stock, c.id, c.name FROM product p 
	JOIN category c ON c.id = p.id_category`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	product := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.ID_Category, &p.Category_Name)
		if err != nil {
			return nil, err
		}

		product = append(product, p)
	}

	return product, nil
}

func (repo *ProductRepository) Create(product *model.Product) error {
	query := "INSERT INTO product (name, price, stock, id_category) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.ID_Category).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) GetByID(id int) (*model.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, c.id, c.name FROM product p
	JOIN category c ON c.id = p.id_category WHERE p.id = $1`

	var p model.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.ID_Category, &p.Category_Name)

	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *model.Product) error {
	query := "UPDATE product SET name = $1, price = $2, stock = $3, id_category = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID_Category, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product not found")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product not found")
	}

	return err
}
