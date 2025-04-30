package services

import (
	"context"
	"time"

	"server/internal/models"
	"server/internal/util"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey = "secret" // should further on store in a separate file
)

type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type LoginUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	AccessToken string // should be private
	ID          string `json:"id"`
	Username    string `json:"username"`
}

type GetUserReq struct {
	ID string `json:"id"`
}

type GetUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

type UserService interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
	GetUserByID(c context.Context, req *GetUserReq) (*GetUserRes, error)
	GetAllUsers(c context.Context) ([]*GetUserRes, error)
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Username:          req.Username,
		EncryptedPassword: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
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

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return &LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.EncryptedPassword) // Вопросики
	if err != nil {
		return &LoginUserRes{}, err
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
		return &LoginUserRes{}, err
	}

	return &LoginUserRes{AccessToken: ss, Username: u.Username, ID: u.ID}, nil
}

func (s *service) GetUserByID(c context.Context, req *GetUserReq) (*GetUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserRes{
		ID:       u.ID,
		Username: u.Username,
		Status:   string(u.Status),
	}, nil
}

func (s *service) GetAllUsers(c context.Context) ([]*GetUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	users, err := s.Repository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*GetUserRes, len(users))
	for _, u := range users {
		result = append(result, &GetUserRes{
			ID:       u.ID,
			Username: u.Username,
			Status:   string(u.Status),
		})
	}

	return result, nil
}
