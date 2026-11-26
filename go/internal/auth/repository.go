package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/mail"
	"strings"
	"time"

	"cinefolio/internal/store"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const Lifetime = 7 * 24 * time.Hour

var ErrUnauthorized = errors.New("invalid credentials or session")

type Repository struct{ db *gorm.DB }

type User struct {
	ID           uint64 `json:"-"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type session struct {
	TokenHash string `gorm:"primaryKey"`
	UserID    uint64
	ExpiresAt int64
}

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) SetPassword(ctx context.Context, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email || len(email) > 254 {
		return store.ValidationError("enter a valid email address")
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	hash, err := hashPassword(ctx, password)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := User{Email: email, PasswordHash: hash}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"password_hash"}),
		}).Create(&user).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ?", user.ID).Delete(&session{}).Error
	})
}

func (r *Repository) Login(ctx context.Context, email, password string) (string, User, error) {
	var user User
	email = strings.ToLower(strings.TrimSpace(email))
	if len(email) > 254 || len(password) == 0 || len(password) > 1024 {
		return "", user, ErrUnauthorized
	}
	db := r.db.WithContext(ctx)
	err := db.Where("email = ?", email).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if _, err := derive(ctx, password, make([]byte, saltBytes)); err != nil {
			return "", user, err
		}
		return "", user, ErrUnauthorized
	}
	if err != nil {
		return "", user, err
	}
	matches, err := verifyPassword(ctx, user.PasswordHash, password)
	if err != nil {
		return "", user, err
	}
	if !matches {
		return "", user, ErrUnauthorized
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", user, err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	now := time.Now()
	err = db.Transaction(func(tx *gorm.DB) error {
		var current User
		if err := tx.Take(&current, user.ID).Error; err != nil {
			return err
		}
		if current.PasswordHash != user.PasswordHash {
			return ErrUnauthorized
		}
		if err := tx.Where("expires_at <= ?", now.Unix()).Delete(&session{}).Error; err != nil {
			return err
		}
		return tx.Create(&session{TokenHash: tokenHash(token), UserID: user.ID, ExpiresAt: now.Add(Lifetime).Unix()}).Error
	})
	return token, user, err
}

func (r *Repository) Current(ctx context.Context, token string) (User, error) {
	var user User
	if len(token) != 43 {
		return user, ErrUnauthorized
	}
	err := r.db.WithContext(ctx).Table("users").
		Select("users.id, users.email").
		Joins("JOIN sessions ON sessions.user_id = users.id").
		Where("sessions.token_hash = ? AND sessions.expires_at > ?", tokenHash(token), time.Now().Unix()).
		Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrUnauthorized
	}
	return user, err
}

func (r *Repository) Logout(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token_hash = ?", tokenHash(token)).Delete(&session{}).Error
}

func tokenHash(token string) string {
	digest := sha256.Sum256([]byte(token))
	return hex.EncodeToString(digest[:])
}
