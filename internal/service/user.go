package service

import (
	"crypto/sha1"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}
type User interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (model.User, error)
	DeleteToken(token string) error
}

type UserService struct {
	repo repository.User
}

func NewUser(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func generateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *UserService) CreateUser(user model.User) (int, error) {
	user.Password = generateHashPassword(user.Password)
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserService) GenerateToken(email string, password string) (string, error) {
	user, err := s.repo.GetUser(email, generateHashPassword(password))
	if err != nil {
		return "", fmt.Errorf("service: generate token: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	user.Token, err = token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return user.Token, nil
}

func (s *UserService) ParseToken(token string) (model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (s *UserService) DeleteToken(token string) error {
	panic("not implemented") // TODO: Implement
}
