package link

import (
	"go/adv-demo/pkg/middleware"
	"go/adv-demo/pkg/req"
	"go/adv-demo/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct{
	LinkRepository *LinkRepository
}

type LinkHandler struct{
	LinkRepository *LinkRepository
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps){
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
	}
	router.HandleFunc("POST /link", handler.Create())
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update()))
}

func (handler *LinkHandler)Create() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		body, err := req.HadleBody[LinkCreateRequest](&w, r)
		if err != nil{
			res.Json(w, err.Error(), 402)
			return
		}

		link := NewLink(body.Url)
		for {
			_, err = handler.LinkRepository.GetByHash(link.Hash)
			if err != nil{
				break
			}
			link.GenereteHash()
		}

		createdLink, err :=handler.LinkRepository.Create(link)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdLink, 201)
	}
}

func (handler *LinkHandler)Delete() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}	

		err = handler.LinkRepository.Delete(uint(id))
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, nil, 200)
	}
}

func (handler *LinkHandler)Update() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		body, err := req.HadleBody[LinkUpdateRequest](&w, r)
		if err != nil{
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{
				ID: uint(id),
			},
			Url: body.Url,
			Hash: body.Hash,
		})
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, link, 200)
	}
}

func (handler *LinkHandler)GoTo() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		hash := r.PathValue("hash")
		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}


