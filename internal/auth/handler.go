package auth

import (
	"go/adv-demo/configs"
	"go/adv-demo/pkg/jwt"
	"go/adv-demo/pkg/req"
	"go/adv-demo/pkg/res"
	"net/http"
)

type AuthHadlerDeps struct{
	*AuthService
	*configs.Config
	*jwt.JWT
}

type AuthHandler struct{
	*AuthService
	*configs.Config
	*jwt.JWT
}

func NewAuthHandler(router *http.ServeMux, deps AuthHadlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
		AuthService: deps.AuthService,
		JWT: deps.JWT,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		body, err := req.HadleBody[LoginRequest](&w, r)
		if err!=nil{
			return
		}

		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		token, err := handler.JWT.Create(jwt.JWTData{
			Email: email,
		})
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := LoginResponse{
			Token: token,
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

		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := handler.JWT.Create(jwt.JWTData{
			Email: email,
		})
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			Token: token,
		}
		res.Json(w, data, 201)
	}
}