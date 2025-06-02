package port

import (
	"Authorization/internal/core/domain"
)

type UserRepository interface {
	GetUserByEmail(email string) (*domain.User, error)
	RegisterUser(regUser domain.User) error
	GetAllUsers() ([]domain.UserReq, error)
	DeleteUser(id int) error
}

type ArticleRepository interface {
	AddArticle(article domain.ArticleReq) error
	GetPublishedArticles() ([]domain.ArticleReq, error)
	GetArticlesInModern() ([]domain.ArticleInModern, error)
	ApproveArticle(id int) error
	RejectArticle(id int, reject domain.RejectReason) error
	GetRejectedArticles(id int) ([]domain.GetReject, error)
	GetEditByID(id int) (domain.EditReject, error)
	EditArticle(edit domain.EditArticleRequest) error
}
