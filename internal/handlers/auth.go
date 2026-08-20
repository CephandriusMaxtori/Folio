package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/CephandriusMaxtori/Folio/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type authResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request")
		return
	}
	if req.Email == "" || req.Password == "" {
		writeError(w, 400, "email and password required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Username:     req.Username,
		Role:         "user",
	}

	count, _ := h.svc.GetStats()
	if count["users"] == 0 {
		user.Role = "admin"
	}

	if err := h.svc.CreateUser(user); err != nil {
		writeError(w, 409, "email already exists")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	writeJSON(w, 201, authResponse{Token: token, User: *user})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request")
		return
	}

	user, err := h.svc.GetUserByEmail(req.Email)
	if err != nil {
		writeError(w, 401, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, 401, "invalid credentials")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	writeJSON(w, 200, authResponse{Token: token, User: *user})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	writeJSON(w, 200, user)
}

func (h *Handler) generateToken(user *models.User) (string, error) {
	secret := h.cfg.JWT.Secret
	if secret == "" {
		secret = os.Getenv("FOLIO_JWT_SECRET")
	}
	if secret == "" {
		secret = "change-me-in-production"
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
