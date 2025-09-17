package models

import (
	"database/sql"
	"errors"
	"rentx/db"
)

type Review struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"userId" binding:"required"`
	PostId   int64  `json:"postId" binding:"required"`
	Review   string `json:"review" binding:"required"`
	DateTime string `json:"dateTime"`
}

// Save inserts a new review
func (r *Review) Save() error {
	res, err := db.DB.Exec(
		"INSERT INTO reviews (userId, postId, review) VALUES (?, ?, ?)",
		r.UserId, r.PostId, r.Review,
	)
	if err != nil {
		return err
	}
	r.Id, _ = res.LastInsertId()
	return nil
}

// Update modifies a review (only by the owner)
func (r *Review) Update(userId int64) error {
	res, err := db.DB.Exec(
		"UPDATE reviews SET review=? WHERE id=? AND userId=?",
		r.Review, r.Id, userId,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("unauthorized or review not found")
	}
	return nil
}

// Delete removes a review (only by the owner)
func (r *Review) Delete(userId int64) error {
	res, err := db.DB.Exec("DELETE FROM reviews WHERE id=? AND userId=?", r.Id, userId)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("unauthorized or review not found")
	}
	return nil
}

// GetReviewByID fetches a single review
func GetReviewByID(id int64) (*Review, error) {
	row := db.DB.QueryRow("SELECT id, userId, postId, review, dateTime FROM reviews WHERE id=?", id)
	var r Review
	if err := row.Scan(&r.Id, &r.UserId, &r.PostId, &r.Review, &r.DateTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("review not found")
		}
		return nil, err
	}
	return &r, nil
}

// ListReviews fetches all reviews for a specific post
func ListReviews(postId int64) ([]Review, error) {
	rows, err := db.DB.Query("SELECT id, userId, postId, review, dateTime FROM reviews WHERE postId=?", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var r Review
		if err := rows.Scan(&r.Id, &r.UserId, &r.PostId, &r.Review, &r.DateTime); err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}
