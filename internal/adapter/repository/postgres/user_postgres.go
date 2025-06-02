package postgres

import (
	"Authorization/internal/core/domain"
	"database/sql"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgreUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func NewPostgreArticleRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.DB.QueryRow(`
		SELECT id, email, password_hash, username, role FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.Role)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) RegisterUser(regUser domain.User) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (email, password_hash, username, role)
		VALUES($1, $2, $3, $4)
	`, regUser.Email, regUser.Password, regUser.Username, regUser.Role)

	return err
}

func (r *PostgresUserRepository) AddArticle(article domain.ArticleReq) error {
	_, err := r.DB.Exec(`
		INSERT INTO articles (title, content, user_id)
		VALUES($1, $2, $3)
	`, article.Title, article.Content, article.ID)

	return err
}

func (r *PostgresUserRepository) GetAllUsers() ([]domain.UserReq, error) {
	rows, err := r.DB.Query(`SELECT id, email, username, role FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.UserReq{}
	for rows.Next() {
		var user domain.UserReq
		if err := rows.Scan(&user.ID ,&user.Email, &user.Username, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *PostgresUserRepository) GetPublishedArticles() ([]domain.ArticleReq, error) {
	rows, err := r.DB.Query(`SELECT a.title, a.content, a.published_at, u.username FROM articles a
							 JOIN users u ON a.user_id = u.id
							 WHERE a.status = 'published' 
							 ORDER BY a.id
							 `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []domain.ArticleReq{}
	for rows.Next() {
		var article domain.ArticleReq
		if err := rows.Scan(&article.Title, &article.Content, &article.Published, &article.Username); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (r *PostgresUserRepository) GetArticlesInModern() ([]domain.ArticleInModern, error) {
	rows, err := r.DB.Query(`SELECT a.id, a.title, a.content, u.username, a.created_at FROM articles a
							JOIN users u ON a.user_id = u.id
							WHERE a.status = 'pending'
							ORDER BY a.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []domain.ArticleInModern{}
	for rows.Next() {
		var a domain.ArticleInModern
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Username, &a.CreatedAt); err != nil {
			return nil, err
		}

		articles = append(articles, a)
	}
	return articles, nil
}

func (r *PostgresUserRepository) GetRejectedArticles(id int) ([]domain.GetReject, error) {
	rows, err := r.DB.Query(`SELECT id, title, content, status, rejection_reason FROM articles WHERE user_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rejects := []domain.GetReject{}
	for rows.Next() {
		var r domain.GetReject
		var rejectionReason sql.NullString

		if err := rows.Scan(&r.ID, &r.Title, &r.Content, &r.Status, &rejectionReason); err != nil {
			return nil, err
		}

		if rejectionReason.Valid {
			r.Reject = rejectionReason.String
		} else {
			r.Reject = ""
		}

		rejects = append(rejects, r)
	}

	return rejects, nil
}

func (r *PostgresUserRepository) GetEditByID(id int) (domain.EditReject ,error) {
	var edit domain.EditReject
	err := r.DB.QueryRow(`SELECT id, title, content, rejection_reason FROM articles WHERE id = $1`, id).Scan(&edit.ID, &edit.Title, &edit.Content, &edit.Reject)
	if err != nil {
		return domain.EditReject{}, err
	}

	return edit, nil
}

func (r *PostgresUserRepository) EditArticle(edit domain.EditArticleRequest) error {
	_, err := r.DB.Exec(`UPDATE articles 
						 SET title = $1, content = $2, status = 'pending', rejection_reason = NULL, rejected_at = NULL, created_at = NOW() 
						 WHERE id = $3`, edit.Title, edit.Content, edit.ID)
	return err
}

func (r *PostgresUserRepository) ApproveArticle(id int) error {
	_, err := r.DB.Exec(`UPDATE articles SET status = 'published', published_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *PostgresUserRepository) RejectArticle(id int, reject domain.RejectReason) error {
	_, err := r.DB.Exec(`UPDATE articles SET status = 'rejected', rejection_reason = $1, rejected_at = NOW() WHERE id = $2`, reject.Reason, id)
	return err
}

func (r *PostgresUserRepository) DeleteUser(id int) error {
	_, err := r.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}