package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"rentx/db"
	"time"
)

type RefreshToken struct {
	Id        int64
	UserId    int64
	Token     string
	ExpiresAt time.Time
	DateTime  string
}

// Create a new random refresh token
func NewRefreshToken(userId int64, daysValid int) (*RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	token := base64.URLEncoding.EncodeToString(b)
	expires := time.Now().Add(time.Hour * 24 * time.Duration(daysValid))
	return &RefreshToken{
		UserId:    userId,
		Token:     token,
		ExpiresAt: expires,
	}, nil
}

// Save refresh token to DB
func (rt *RefreshToken) Save() error {
	query := "INSERT INTO refreshTokens (userId, token, expiresAt) VALUES (?, ?, ?)"
	res, err := db.DB.Exec(query, rt.UserId, rt.Token, rt.ExpiresAt)
	if err != nil {
		return err
	}
	rt.Id, _ = res.LastInsertId()
	return nil
}

// Get a refresh token from DB and verify expiration
func GetRefreshToken(token string) (*RefreshToken, error) {
	row := db.DB.QueryRow("SELECT id, userId, token, expiresAt, dateTime FROM refreshTokens WHERE token=?", token)
	var rt RefreshToken
	if err := row.Scan(&rt.Id, &rt.UserId, &rt.Token, &rt.ExpiresAt, &rt.DateTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}
	return &rt, nil
}
