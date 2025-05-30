package app

import (
	"Authorization/internal/adapter/handler"
	"net/http"
)

func setupRoutex(mux *http.ServeMux, handler *handler.Handler) {
	mux.HandleFunc("GET /me", handler.GetProfile)
	mux.HandleFunc("POST /login", handler.LoginHandler)
	mux.HandleFunc("GET /articles", handler.GetArticle)

	mux.HandleFunc("POST /admin/users", handler.AddUser)
	mux.HandleFunc("GET /admin/users", handler.GetAllUsers)

	mux.HandleFunc("GET /moderator/articles", handler.GetArticleModern)
	mux.HandleFunc("POST /moderator/articles/{id}/approve", handler.PostArtilceModernApprove)
	mux.HandleFunc("POST /moderator/articles/{id}/reject", handler.PostArtilceModernReject)

	mux.HandleFunc("POST /author/articles", handler.PostAuthorArticle)
	mux.HandleFunc("GET /author/articles/all", handler.GetRejectedArticles)
	mux.HandleFunc("GET /author/articles/{id}", handler.GetEditArticle)
	mux.HandleFunc("POST /author/articles/{id}/edit", handler.PostEditedArticle)

	mux.HandleFunc("POST /logout", handler.LogOut)

	mux.Handle("GET /", http.FileServer(http.Dir("./front")))
}
