package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

type MockJWT struct {
	token string
	err   error
 }

func (m *MockJWT) GenerateToken(userID int) (string, error) {
	return "mockToken", nil
}

func (m *MockJWT) ParseToken(tokenString string) (int, error) {
	return 1, nil
}

type MockStorage struct {
	userID int
	err    error
 }
 
func (m MockStorage) Login(ctx context.Context, login, password string) (int, error) {
	return m.userID, m.err
}

func (m MockStorage) CheckLogin(ctx context.Context, login string) error {
	return nil
}

func (m *MockStorage) Register(ctx context.Context, login, password string) (int, error) {
	return 1, nil
}

func TestGetHashedPassword(t *testing.T) {
	user := NewUserAuth(&MockStorage{}, &MockJWT{})

	password := "mySecurePassword"
	hashedPassword := user.GetHashedPassword(password)

	hash := sha256.New()
	hash.Write([]byte(password))
	expectedHash := hex.EncodeToString(hash.Sum(nil))

	if hashedPassword != expectedHash {
		t.Errorf("expected %s, got %s", expectedHash, hashedPassword)
	}
}

func TestLoginUser(t *testing.T) {
       tests := []struct {
           username string
           password string
           userID   int
           tokenErr error
           wantErr  bool
       }
       {
           username: "testuser",
           password: "password",
           userID:   1,
           tokenErr: nil,
           wantErr:  false,
       },
       {
           username: "wronguser",
           password: "wrongpassword",
           userID:   0,
           tokenErr: nil,
           wantErr:  true,
       },
       {
           username: "testuser",
           password: "password",
           userID:   1,
           tokenErr: errors.New("token error"),
           wantErr:  true,
       },
       }

       for _, tt := range tests {
           storage := &MockStorage{userID: tt.userID, err: tt.wantErr}
           jwt := &MockJWT{token: "mockToken", err: tt.tokenErr}
           user := NewUserAuth(storage, jwt)

           token, err := user.LoginUser(context.Background(), tt.username, tt.password)

           if (err != nil) != tt.wantErr {
               t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
           } else if !tt.wantErr && token != "mockToken" {
               t.Errorf("LoginUser() token = %v, want %v", token, "mockToken")
           }
       } }