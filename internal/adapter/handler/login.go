package handler

import (
	"Authorization/internal/core/domain"
	"Authorization/internal/core/port"
	"Authorization/internal/core/util"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	UserRepo    port.UserRepository
	Logger      *slog.Logger
	ArticleRepo port.ArticleRepository
}

func NewHandler(repo port.UserRepository, log *slog.Logger, article port.ArticleRepository) *Handler {
	return &Handler{
		UserRepo:    repo,
		Logger:      log,
		ArticleRepo: article,
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if strings.ContainsAny(req.Email, "'\"--") {
		h.Logger.Info("Подозоительный email", "email", req.Email)
		http.Error(w, "Suspence email", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	isMatch := util.CheckPasswordHash(req.Password, user.Password)
	if !isMatch {
		h.Logger.Info("Кто то пытался зайти в аккаунт. Пароль был не правильным", "Аккаунт", req.Email)
		http.Error(w, "Invalid password", http.StatusForbidden)
		return
	}

	token := util.GenerateToken(user)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	h.Logger.Info("Пользователь зашел в аккаунт", "email", req.Email)
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{"message": "OK"}`))
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	articles, err := h.ArticleRepo.GetPublishedArticles()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"articles": articles}); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed encode articles", http.StatusInternalServerError)
	}
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil || claims["role"] != "admin" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	var req domain.AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Role == "" || req.Username == "" {
		http.Error(w, "Fill in the entire field", http.StatusBadRequest)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed hash password", http.StatusInternalServerError)
		return
	}

	user := domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Username: req.Username,
		Role:     req.Role,
	}

	err = h.UserRepo.RegisterUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Админ создал пользователя", "email", user.Email, "username", user.Username, "role", user.Role)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"email":    claims["email"],
		"role":     claims["role"],
		"username": claims["username"],
	})
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil || claims["role"] != "admin" {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}

	users, err := h.UserRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}

func (h *Handler) GetArticleModern(w http.ResponseWriter, r *http.Request) {
	articles, err := h.ArticleRepo.GetArticlesInModern()
	if err != nil {
		http.Error(w, "Failed to load artic in modern", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"articles": articles})
}

func (h *Handler) PostArtilceModernApprove(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil || claims["role"] != "moderator" {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}

	if err := h.ArticleRepo.ApproveArticle(id); err != nil {
		http.Error(w, "Failed approve article", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Модератор проверил и одобрил статью", "модератор", claims["email"])
}

func (h *Handler) PostArtilceModernReject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil || claims["role"] != "moderator" {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}

	var reject domain.RejectReason
	if err := json.NewDecoder(r.Body).Decode(&reject); err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	if err := h.ArticleRepo.RejectArticle(id, reject); err != nil {
		http.Error(w, "Failed reject article", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Модератор проверил и не одобрил статью", "модератор", claims["email"])
}

func (h *Handler) PostAuthorArticle(w http.ResponseWriter, r *http.Request) {
	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil || claims["role"] != "author" {
		http.Error(w, "Invavlid role", http.StatusForbidden)
		return
	}

	var article domain.ArticleReq
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	isDangerTitle := util.CheckInput(article.Title)
	if isDangerTitle {
		h.Logger.Info("Был откланен опасный тайтл", "отправитель", claims["email"])
		http.Error(w, "Dangerous message", http.StatusBadRequest)
		return
	}

	isDangerContent := util.CheckInput(article.Content)
	if isDangerContent {
		h.Logger.Info("Был откланен опасный контент", "отправитель", claims["email"])
		http.Error(w, "Dangerous message", http.StatusBadRequest)
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	article.ID = int(userID)

	if err := h.ArticleRepo.AddArticle(article); err != nil {
		http.Error(w, "Failed add article", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Автор отправил свою работу на проверку", "автор", claims["email"])
}

func (h *Handler) GetRejectedArticles(w http.ResponseWriter, r *http.Request) {
	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil || claims["role"] != "author" {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}

	id := int(claims["user_id"].(float64))

	reject, err := h.ArticleRepo.GetRejectedArticles(id)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Failed to get rejected articles", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"articles": reject})
}

func (h *Handler) GetEditArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil || claims["role"] != "author" {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}

	edit, err := h.ArticleRepo.GetEditByID(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get edited article", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"include": edit})
}

func (h *Handler) PostEditedArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil || claims["role"] != "author" {
		http.Error(w, "Invavlid role", http.StatusForbidden)
		return
	}

	var edit domain.EditArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&edit); err != nil {
		http.Error(w, "Failed decode edit", http.StatusBadRequest)
		return
	}

	isDangerTitle := util.CheckInput(edit.Title)
	if isDangerTitle {
		h.Logger.Info("Был откланен опасный тайтл", "отправитель", claims["email"])
		http.Error(w, "Dangerous message", http.StatusBadRequest)
		return
	}

	isDangerContent := util.CheckInput(edit.Content)
	if isDangerContent {
		h.Logger.Info("Был откланен опасный контент", "отправитель", claims["email"])
		http.Error(w, "Dangerous message", http.StatusBadRequest)
		return
	}

	edit.ID = id
	if err := h.ArticleRepo.EditArticle(edit); err != nil {
		http.Error(w, "Failed edit article", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Автор исправил статью которую не одобрили", "автор", claims["email"])
}

func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	claims, err := util.ValidateTokenFromRequest(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}

	h.Logger.Info("Пользователь вышел из аккаунта", "email", claims["email"])

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}
