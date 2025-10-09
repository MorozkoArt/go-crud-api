package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/MorozkoArt/go-crud-api/internal/auth"
)

type userKey string

const (
    UserIDKey    userKey = "user_id"
    UserEmailKey userKey = "user_email"
)

func AuthMiddleware(jwtAuth *auth.JWTService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, `{"success": false, "error": "Authorization header required"}`, http.StatusUnauthorized)
                return
            }
            
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, `{"success": false, "error": "Invalid authorization format"}`, http.StatusUnauthorized)
                return
            }
            
            tokenString := parts[1]
            
            claims, err := jwtAuth.ValidateToken(tokenString)
            if err != nil {
                http.Error(w, `{"success": false, "error": "Invalid or expired token"}`, http.StatusUnauthorized)
                return
            }
            
            ctx := r.Context()
            ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
            ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserID(ctx context.Context) (int64, bool) {
    userID, ok := ctx.Value(UserIDKey).(int64)
    return userID, ok
}

func GetUserEmail(ctx context.Context) (string, bool) {
    email, ok := ctx.Value(UserEmailKey).(string)
    return email, ok
}