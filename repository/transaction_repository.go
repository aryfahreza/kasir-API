package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"kasir-api/model"
	"strings"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []model.CheckoutItem) (*model.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]model.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM product WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product ID %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE product SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transaction (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	var (
		query  = "INSERT INTO transaction_detail (transaction_id, product_id, quantity, subtotal) VALUES "
		args   []interface{}
		values []string
	)

	for i, d := range details {
		base := i * 4

		values = append(values,
			fmt.Sprintf("($%d,$%d,$%d,$%d)",
				base+1, base+2, base+3, base+4,
			),
		)

		args = append(args,
			transactionID,
			d.ProductID,
			d.Quantity,
			d.Subtotal,
		)

		query += strings.Join(values, ",")
		query += " RETURNING id"
	}

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		details[i].ID = id
		details[i].TransactionID = transactionID
		i++
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetTransactionReport() (*model.TransactionReport, error) {
	query := `SELECT SUM(t.total_amount) AS total_revenue, COUNT(t) AS total_transaction, p.name, SUM(td.quantity) AS sold_qty
				FROM transaction_detail td
				JOIN transaction t ON t.id = td.transaction_id
				JOIN product p ON p.id = td.product_id
				GROUP BY p.name
				ORDER BY COUNT(td.product_id) DESC
				LIMIT 1`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	total_revenue := 0
	total_transaction := 0
	var b model.BestSelling
	for rows.Next() {
		err := rows.Scan(
			&total_revenue,
			&total_transaction,
			&b.Name,
			&b.Quantity,
		)

		if err != nil {
			return nil, err
		}
	}

	var result model.TransactionReport
	result.TotalRevenue = total_revenue
	result.TotalTransaction = total_transaction
	result.BestSelling = b

	return &result, err
}

func (repo *TransactionRepository) GetTransactionReportByDate(startDate string, endDate string) (*model.TransactionReport, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	end = end.AddDate(0, 0, 1)

	query := `SELECT SUM(t.total_amount) AS total_revenue, COUNT(t) AS total_transaction, p.name, SUM(td.quantity) AS sold_qty
				FROM transaction_detail td
				JOIN transaction t ON t.id = td.transaction_id
				JOIN product p ON p.id = td.product_id
				WHERE t.created_at >= $1
				AND t.created_at < $2
				GROUP BY p.name
				ORDER BY COUNT(td.product_id) DESC
				LIMIT 1`

	rows, err := repo.db.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	total_revenue := 0
	total_transaction := 0
	var b model.BestSelling
	var found bool
	for rows.Next() {
		found = true

		err := rows.Scan(
			&total_revenue,
			&total_transaction,
			&b.Name,
			&b.Quantity,
		)

		if err != nil {
			return nil, err
		}
	}

	if !found {
		return nil, errors.New("Report not found, please select another date")
	}

	var result model.TransactionReport
	result.TotalRevenue = total_revenue
	result.TotalTransaction = total_transaction
	result.BestSelling = b

	return &result, err
}
