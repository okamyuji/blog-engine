package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"my-blog-engine/internal/interface/middleware"
	"my-blog-engine/internal/interface/presenter"
	"my-blog-engine/internal/usecase"
)

// AuthHandler 認証ハンドラー
type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

// NewAuthHandler 新しいAuthHandlerを作成
func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshTokenRequest リフレッシュトークンリクエスト
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// Login ログインハンドラー
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 入力検証
	if req.Username == "" || req.Password == "" {
		presenter.JSONError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// ログイン処理
	response, err := h.authUseCase.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			presenter.JSONError(w, http.StatusUnauthorized, "Invalid username or password")
		} else if strings.Contains(err.Error(), "inactive") {
			presenter.JSONError(w, http.StatusForbidden, "User account is inactive")
		} else {
			presenter.JSONError(w, http.StatusInternalServerError, "Login failed")
		}
		return
	}

	presenter.JSONResponse(w, http.StatusOK, response)
}

// Logout ログアウトハンドラー
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Authorizationヘッダーからトークンを取得
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		presenter.JSONError(w, http.StatusUnauthorized, "Authorization header required")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		presenter.JSONError(w, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}

	token := parts[1]

	// ログアウト処理
	if err := h.authUseCase.Logout(r.Context(), token); err != nil {
		presenter.JSONError(w, http.StatusInternalServerError, "Logout failed")
		return
	}

	presenter.JSONSuccess(w, nil, "Logged out successfully")
}

// RefreshToken トークンリフレッシュハンドラー
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenter.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		presenter.JSONError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// トークンリフレッシュ処理
	response, err := h.authUseCase.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "revoked") {
			presenter.JSONError(w, http.StatusUnauthorized, "Invalid or revoked refresh token")
		} else {
			presenter.JSONError(w, http.StatusInternalServerError, "Token refresh failed")
		}
		return
	}

	presenter.JSONResponse(w, http.StatusOK, response)
}

// Me 現在のユーザー情報を取得
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		presenter.JSONError(w, http.StatusUnauthorized, "User not found")
		return
	}

	presenter.JSONResponse(w, http.StatusOK, user)
}
