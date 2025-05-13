package services

import (
	"context"
	"fmt"
	"time"

	"chatgo/server/internal/interfaces"
	"chatgo/server/internal/models"
	"chatgo/server/internal/util"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey = "secret" // should further on store in a separate file
)

func (s *service) CreateUser(c context.Context, req *interfaces.CreateUserReq) (*interfaces.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Username:          req.Username,
		EncryptedPassword: hashedPassword,
		Status:            models.UserStatus(models.Offline),
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &interfaces.CreateUserRes{
		ID:       r.ID,
		Username: r.Username,
	}

	return res, nil
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *interfaces.LoginUserReq) (*interfaces.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		fmt.Println("After getting user by username")
		return &interfaces.LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.EncryptedPassword)
	if err != nil {
		fmt.Println("After checking password")
		return &interfaces.LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("After signing token")
		return &interfaces.LoginUserRes{}, err
	}

	return &interfaces.LoginUserRes{AccessToken: ss, Username: u.Username, ID: u.ID}, nil
}

func (s *service) GetUserByID(c context.Context, req *interfaces.GetUserReq) (*interfaces.GetUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &interfaces.GetUserRes{
		ID:       u.ID,
		Username: u.Username,
	}, nil
}

func (s *service) GetAllUsers(c context.Context) ([]*interfaces.GetUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	users, err := s.Repository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*interfaces.GetUserRes, 0, len(users))
	for _, u := range users {
		result = append(result, &interfaces.GetUserRes{
			ID:       u.ID,
			Username: u.Username,
		})
	}

	return result, nil
}
