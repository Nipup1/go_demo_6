package main

import (
	"fmt"
	"go/adv-demo/configs"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/link"
	"go/adv-demo/internal/stat"
	"go/adv-demo/internal/user"
	"go/adv-demo/pkg/db"
	"go/adv-demo/pkg/event"
	"go/adv-demo/pkg/jwt"
	"go/adv-demo/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	jwt := jwt.NewJWT(conf.Auth.Secret)
	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	//Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	//Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		StatRepository: statRepository,
		EventBus: eventBus,
	})

	//Handler
	auth.NewAuthHandler(router, auth.AuthHadlerDeps{
		Config: conf,
		AuthService: authService,
		JWT: jwt,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config: conf,
		EventBus: eventBus,
	})
	stat.NewStatHendler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config: conf,
	})

	//Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	//Events
	go statService.AddClick()

	server := http.Server{
		Addr: ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
