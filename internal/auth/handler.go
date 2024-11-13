package auth

import (
	"fmt"
	"go/adv-demo/configs"
	"go/adv-demo/pkg/req"
	"go/adv-demo/pkg/res"
	"net/http"
)

type AuthHadlerDeps struct{
	*AuthService
	*configs.Config
}

type AuthHandler struct{
	*AuthService
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHadlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		payload, err := req.HadleBody[LoginRequest](&w, r)
		if err!=nil{
			return
		}

		fmt.Println(payload)

		data := LoginResponse{
			Token: "123",
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		body, err := req.HadleBody[RegisterRequest](&w, r)
		if err!=nil{
			return
		}

		handler.AuthService.Register(body.Email, body.Password, body.Name)
	}
}