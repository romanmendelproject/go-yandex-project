package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	log "github.com/sirupsen/logrus"
)

type JWT interface {
	GenerateToken(userID int) (string, error)
	ParseToken(tokenString string) (int, error)
}

type Storage interface {
	Login(ctx context.Context, login, password string) (int, error)
	CheckLogin(ctx context.Context, login string) error
	Register(ctx context.Context, login, password string) (int, error)
}

type User struct {
	db    Storage
	token JWT
}

func NewUserAuth(db Storage, token JWT) *User {
	return &User{
		db:    db,
		token: token,
	}
}

// GetHashedPassword - generate hash from password to store secure data
func (u *User) GetHashedPassword(password string) string {
	hash := sha256.New()

	hash.Write([]byte(password))

	return hex.EncodeToString(hash.Sum(nil))
}

func (u *User) LoginUser(ctx context.Context, username, password string) (string, error) {
	// generating hash from password
	hashPassword := u.GetHashedPassword(password)

	userID, err := u.db.Login(ctx, username, hashPassword)
	if err != nil {
		return "", err
	}

	// if userID == 0 => user not exists
	if userID == 0 {
		log.Error("user not found", "username", username)

		return "", errors.New("user not found")
	}

	// generating token
	return u.token.GenerateToken(userID)
}

func (u *User) RegisterUser(ctx context.Context, username string, password string) (string, error) {
	// check if login already occupied
	if err := u.db.CheckLogin(ctx, username); err != nil {
		return "", err
	}

	// generating hash from password
	hashPassword := u.GetHashedPassword(password)

	// registering user
	id, err := u.db.Register(ctx, username, hashPassword)
	if err != nil {
		return "", err
	}

	return u.token.GenerateToken(id)
}
