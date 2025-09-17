package models

import (
	"errors"
	"rentx/db"
)

type Category struct {
	Id     int64  `json:"id"`
	UserId int64  `json:"-"`
	Name   string `json:"name" binding:"required"`
}

// Authorization error
var ErrUnauthorized = errors.New("unauthorized action")

// ---------- Category Methods ----------

func (c *Category) Save() error {
	query := "INSERT INTO categories (name) VALUES (?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(c.Name)
	if err != nil {
		return err
	}
	c.Id, _ = res.LastInsertId()
	return nil
}

func (c *Category) Update(userId int64) error {
	if c.UserId != userId {
		return ErrUnauthorized
	}
	query := "UPDATE categories SET name = ? WHERE id = ?"
	_, err := db.DB.Exec(query, c.Name, c.Id)
	return err
}

func (c *Category) Delete(userId int64) error {
	if c.UserId != userId {
		return ErrUnauthorized
	}
	query := "DELETE FROM categories WHERE id = ?"
	_, err := db.DB.Exec(query, c.Id)
	return err
}

func GetCategories() ([]Category, error) {
	rows, err := db.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
