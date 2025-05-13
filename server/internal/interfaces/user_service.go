package interfaces

import "context"

type UserService interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
	GetUserByID(c context.Context, req *GetUserReq) (*GetUserRes, error)
	GetAllUsers(c context.Context) ([]*GetUserRes, error)
}
