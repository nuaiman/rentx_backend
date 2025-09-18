package models

import (
	"database/sql"
	"errors"
	"rentx/db"
	"rentx/utils"
)

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
	Image    string `json:"image"`
	Role     string `json:"role"` // 'user' | 'admin' | 'superadmin'
	DateTime string `json:"dateTime"`
}

// Save inserts a new user into the database
func (u *User) Save() error {
	if u.Role == "" {
		u.Role = "user"
	}
	if u.Image == "" {
		u.Image = ""
	}
	hashedPassword, err := utils.GenerateHashword(u.Password)
	if err != nil {
		return err
	}
	query := "INSERT INTO users (name, email, phone, password, image, role) VALUES (?, ?, ?, ?, ?, ?)"
	res, err := db.DB.Exec(query, u.Name, u.Email, u.Phone, hashedPassword, u.Image, u.Role)
	if err != nil {
		return err
	}
	u.Id, _ = res.LastInsertId()
	return nil
}

func (u *User) LoadByEmail() error {
	query := "SELECT id, name, email, phone, password, image, role FROM users WHERE email = ?"
	return db.DB.QueryRow(query, u.Email).Scan(
		&u.Id, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Image, &u.Role,
	)
}

// GetUserByID fetches a user by ID
func GetUserByID(id int64) (*User, error) {
	row := db.DB.QueryRow("SELECT id, name, email, phone, password, image, role, dateTime FROM users WHERE id=?", id)
	var u User
	if err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Image, &u.Role, &u.DateTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

// ListUsers returns all users
func ListUsers() ([]User, error) {
	rows, err := db.DB.Query("SELECT id, name, email, phone, image, role, dateTime FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Phone, &u.Image, &u.Role, &u.DateTime); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// ---------------- Role Management & Permissions ---------------- //

// PromoteUserToAdmin promotes a user to admin (only superadmin can call)
func (user *User) PromoteUserToAdmin(targetId int64) error {
	if user.Role != "superadmin" {
		return errors.New("only superadmin can promote to admin")
	}

	target, err := GetUserByID(targetId)
	if err != nil {
		return err
	}

	if target.Role == "superadmin" {
		return errors.New("cannot change role of another superadmin")
	}

	_, err = db.DB.Exec("UPDATE users SET role='admin' WHERE id=?", targetId)
	return err
}

// PromoteUserToSuperadmin promotes a user to superadmin (only superadmin can call)
func (user *User) PromoteUserToSuperadmin(targetId int64) error {
	if user.Role != "superadmin" {
		return errors.New("only superadmin can promote to superadmin")
	}

	_, err := db.DB.Exec("UPDATE users SET role='superadmin' WHERE id=?", targetId)
	return err
}

// DeleteUser deletes a user (role-based permissions)
func (user *User) DeleteUser(targetId int64) error {
	target, err := GetUserByID(targetId)
	if err != nil {
		return err
	}

	switch user.Role {
	case "superadmin":
		// superadmin can delete anyone
		_, err := db.DB.Exec("DELETE FROM users WHERE id=?", targetId)
		return err
	case "admin":
		if target.Role == "superadmin" {
			return errors.New("admin cannot delete superadmin")
		}
		_, err := db.DB.Exec("DELETE FROM users WHERE id=?", targetId)
		return err
	default:
		return errors.New("user does not have permission to delete")
	}
}

// Helper methods
func (u *User) IsAdminOrSuperadmin() bool {
	return u.Role == "admin" || u.Role == "superadmin"
}

func (u *User) IsSuperadmin() bool {
	return u.Role == "superadmin"
}
