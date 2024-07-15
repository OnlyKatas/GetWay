package controller

import (
	"getway/internal/modules/controller/auth"
	"getway/internal/modules/controller/pay"
	"getway/internal/modules/controller/users"
)

type Handler struct {
	Auth *auth.AuthController
	User *users.UserController
	Pay  *pay.Controller
}

func NewHandler(as auth.AuthServicer, us users.UserServicer, ps pay.PayServicer) *Handler {
	return &Handler{
		Auth: auth.NewAuthConroller(as, ps),
		User: users.NewUsersController(us),
		Pay:  pay.NewController(ps),
	}
}
