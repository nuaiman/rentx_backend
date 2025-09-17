package models

import (
	"database/sql"
	"errors"
	"rentx/db"
)

// Order represents a single order
type Order struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"userId" binding:"required"`
	PostId   int64  `json:"postId" binding:"required"`
	DateTime string `json:"dateTime"`
}

// Create inserts a new order into the database
func (o *Order) Create() error {
	res, err := db.DB.Exec(
		"INSERT INTO orders (userId, postId) VALUES (?, ?)",
		o.UserId, o.PostId,
	)
	if err != nil {
		return err
	}
	o.Id, _ = res.LastInsertId()
	return nil
}

// Delete removes an order from the database
func (o *Order) Delete() error {
	_, err := db.DB.Exec("DELETE FROM orders WHERE id=?", o.Id)
	return err
}

// GetOrder fetches a single order by ID
func GetOrder(id int64) (*Order, error) {
	row := db.DB.QueryRow("SELECT id, userId, postId, dateTime FROM orders WHERE id=?", id)
	var o Order
	if err := row.Scan(&o.Id, &o.UserId, &o.PostId, &o.DateTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &o, nil
}

// ListOrders fetches all orders
func ListOrders() ([]Order, error) {
	rows, err := db.DB.Query("SELECT id, userId, postId, dateTime FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.Id, &o.UserId, &o.PostId, &o.DateTime); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
