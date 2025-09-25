package models

import (
	"database/sql"
	"errors"
	"rentx/db"
)

type Post struct {
	Id           int64    `json:"id"`
	UserId       int64    `json:"userId"`
	CategoryId   int64    `json:"categoryId" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Address      string   `json:"address" binding:"required"`
	Description  string   `json:"description" binding:"required"`
	DailyPrice   float64  `json:"dailyPrice"`
	WeeklyPrice  float64  `json:"weeklyPrice"`
	MonthlyPrice float64  `json:"monthlyPrice"`
	ImageUrls    []string `json:"imageUrls"`
	Status       string   `json:"status"`
	DateTime     string   `json:"dateTime"`
}

// Save inserts a new post with images
func (p *Post) Save(role string) error {
	// Ensure new posts have 'pending' status by default for normal users. Else, auto approvede
	if role == "admin" || role == "superadmin" {
		p.Status = "approved"
	} else {
		p.Status = "pending"
	}
	// Insert post
	res, err := db.DB.Exec(`
		INSERT INTO posts 
		(userId, categoryId, name, address, description, dailyPrice, weeklyPrice, monthlyPrice, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.UserId, p.CategoryId, p.Name, p.Address, p.Description, p.DailyPrice, p.WeeklyPrice, p.MonthlyPrice, p.Status)
	if err != nil {
		return err
	}
	p.Id, _ = res.LastInsertId()

	// Insert images if any
	if len(p.ImageUrls) > 0 {
		tx, err := db.DB.Begin()
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare(`
			INSERT INTO post_images (postId, imageUrl, position) VALUES (?, ?, ?)`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		for i, url := range p.ImageUrls {
			if url == "" {
				continue
			}
			if _, err := stmt.Exec(p.Id, url, i); err != nil {
				tx.Rollback()
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// Update updates a post (owner or admin)
func (p *Post) Update(userId int64, role string) error {
	query := `
		UPDATE posts SET categoryId=?, name=?, address=?, description=?, dailyPrice=?, weeklyPrice=?, monthlyPrice=?, status=?
		WHERE id=?`
	args := []interface{}{p.CategoryId, p.Name, p.Address, p.Description, p.DailyPrice, p.WeeklyPrice, p.MonthlyPrice, p.Status, p.Id}

	if role != "admin" && role != "superadmin" {
		query += " AND userId=?"
		args = append(args, userId)
	}

	res, err := db.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		return errors.New("unauthorized or post not found")
	}

	_, err = db.DB.Exec("DELETE FROM post_images WHERE postId=?", p.Id)
	if err != nil {
		return err
	}

	for i, url := range p.ImageUrls {
		if url == "" {
			continue
		}
		_, err := db.DB.Exec(`
			INSERT INTO post_images (postId, imageUrl, position)
			VALUES (?, ?, ?)`, p.Id, url, i)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete removes a post (owner or admin)
func (p *Post) Delete(userId int64, role string) error {
	query := "DELETE FROM posts WHERE id=?"
	args := []interface{}{p.Id}

	// Only restrict to owner if not admin/superadmin
	if role != "admin" && role != "superadmin" {
		query += " AND userId=?"
		args = append(args, userId)
	}

	res, err := db.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("unauthorized or post not found")
	}
	return nil
}

// GetPostByID fetches a single post with images
func GetPostByID(id int64) (*Post, error) {
	row := db.DB.QueryRow(`
		SELECT id, userId, categoryId, name, address, description, dailyPrice, weeklyPrice, monthlyPrice, status, dateTime
		FROM posts WHERE id=? AND status='approved'`, id)
	var p Post
	err := row.Scan(&p.Id, &p.UserId, &p.CategoryId, &p.Name, &p.Address, &p.Description, &p.DailyPrice, &p.WeeklyPrice, &p.MonthlyPrice, &p.Status, &p.DateTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	rows, err := db.DB.Query(`SELECT imageUrl FROM post_images WHERE postId=? ORDER BY position ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p.ImageUrls = []string{}
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		p.ImageUrls = append(p.ImageUrls, url)
	}

	return &p, nil
}

// ListApprovedPosts fetches all approved posts with images
func ListApprovedPosts() ([]Post, error) {
	rows, err := db.DB.Query(`
		SELECT id, userId, categoryId, name, address, description, dailyPrice, weeklyPrice, monthlyPrice, status, dateTime 
		FROM posts WHERE status='approved'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.UserId, &p.CategoryId, &p.Name, &p.Address, &p.Description, &p.DailyPrice, &p.WeeklyPrice, &p.MonthlyPrice, &p.Status, &p.DateTime); err != nil {
			return nil, err
		}

		imageRows, err := db.DB.Query(`SELECT imageUrl FROM post_images WHERE postId=? ORDER BY position ASC`, p.Id)
		if err != nil {
			return nil, err
		}
		defer imageRows.Close()

		p.ImageUrls = []string{}
		for imageRows.Next() {
			var url string
			if err := imageRows.Scan(&url); err != nil {
				return nil, err
			}
			p.ImageUrls = append(p.ImageUrls, url)
		}
		defer imageRows.Close()
		posts = append(posts, p)
	}

	return posts, nil
}

// ListPendingPosts returns all posts with status "pending"
func ListPendingPosts() ([]Post, error) {
	rows, err := db.DB.Query(`
        SELECT id, userId, categoryId, name, address, description, dailyPrice, weeklyPrice, monthlyPrice, status, dateTime 
        FROM posts WHERE status='pending'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.UserId, &p.CategoryId, &p.Name, &p.Address, &p.Description,
			&p.DailyPrice, &p.WeeklyPrice, &p.MonthlyPrice, &p.Status, &p.DateTime); err != nil {
			return nil, err
		}

		// fetch images
		imageRows, err := db.DB.Query(`SELECT imageUrl FROM post_images WHERE postId=? ORDER BY position ASC`, p.Id)
		if err != nil {
			return nil, err
		}
		defer imageRows.Close()

		p.ImageUrls = []string{}
		for imageRows.Next() {
			var url string
			if err := imageRows.Scan(&url); err != nil {
				return nil, err
			}
			p.ImageUrls = append(p.ImageUrls, url)
		}

		posts = append(posts, p)
	}

	return posts, nil
}

// UpdateStatus updates the status of a post (approved/rejected)
func UpdateStatus(postID int64, status string) error {
	res, err := db.DB.Exec(`UPDATE posts SET status=? WHERE id=?`, status, postID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}

// GetPostStatus returns the current status of a post
func GetPostStatus(postID int64) (string, error) {
	var status string
	err := db.DB.QueryRow(`SELECT status FROM posts WHERE id=?`, postID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("post not found")
		}
		return "", err
	}
	return status, nil
}

// ListPendingPosts returns all posts with status "pending"
// func ListAllPosts() ([]Post, error) {
// 	rows, err := db.DB.Query(`
//         SELECT id, userId, categoryId, name, address, description, dailyPrice, weeklyPrice, monthlyPrice, status, dateTime
//         FROM posts WHERE status='pending'`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var posts []Post
// 	for rows.Next() {
// 		var p Post
// 		if err := rows.Scan(&p.Id, &p.UserId, &p.CategoryId, &p.Name, &p.Address, &p.Description,
// 			&p.DailyPrice, &p.WeeklyPrice, &p.MonthlyPrice, &p.Status, &p.DateTime); err != nil {
// 			return nil, err
// 		}
// 		// fetch images
// 		imageRows, err := db.DB.Query(`SELECT imageUrl FROM post_images WHERE postId=? ORDER BY position ASC`, p.Id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer imageRows.Close()
// 		p.ImageUrls = []string{}
// 		for imageRows.Next() {
// 			var url string
// 			if err := imageRows.Scan(&url); err != nil {
// 				return nil, err
// 			}
// 			p.ImageUrls = append(p.ImageUrls, url)
// 		}
// 		posts = append(posts, p)
// 	}
// 	return posts, nil
// }
