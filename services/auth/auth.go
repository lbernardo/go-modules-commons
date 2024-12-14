package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	secret string
	logger *zap.Logger
}

type Claims struct {
	Sub   string                 `json:"sub"`
	Email string                 `json:"email"`
	Data  map[string]interface{} `json:"data"`
	Exp   any                    `json:"exp"`
}

func NewAuth(cfg *viper.Viper, logger *zap.Logger) *Auth {
	return &Auth{secret: cfg.GetString("app.auth.secret"), logger: logger}
}

// GenerateToken generate a new token with claims and return this
func (a *Auth) GenerateToken(id string, username string, data map[string]any, options ...GenerateTokenOption) (string, error) {
	mapClaims := jwt.MapClaims{
		"sub":      id,
		"email":    username,
		"username": username,
		"data":     data,
	}
	defaultOptions := &generateTokenOption{
		DoesNotExpire: false,
		Expiration:    time.Now().Add(time.Hour * 24),
	}

	for _, option := range options {
		option(defaultOptions)
	}
	if !defaultOptions.DoesNotExpire {
		mapClaims["exp"] = defaultOptions.Expiration.Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenString, err := token.SignedString([]byte(a.secret))
	return tokenString, err
}

// ValidateToken validate if token is ok and return jwt Token if success
func (a *Auth) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		a.logger.Error("token parse failed", zap.Error(err))
		return nil, err
	}
	if !token.Valid {
		a.logger.Error("token is invalid")
		return nil, errors.New("invalid-token")
	}
	return token, nil
}

// GetClaims return claims
func (a *Auth) GetClaims(tokenString string) (*Claims, error) {
	token, err := a.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims-not-found")
	}

	return &Claims{
		Sub:   claims["sub"].(string),
		Email: claims["email"].(string),
		Data:  claims["data"].(map[string]interface{}),
		Exp:   claims["exp"],
	}, nil
}

func (a *Auth) checkMiddleware(r *gin.Context) (*Claims, error) {
	token := r.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil, errors.New("token-not-found")
	}
	return a.GetClaims(token)
}

// MiddlewareDefault is the middleware to use in endpoints with authorization
func (a *Auth) MiddlewareDefault(r *gin.Context) {
	claims, err := a.checkMiddleware(r)
	if err != nil {
		r.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	r.Set("id", claims.Sub)
	r.Set("email", claims.Email)
	r.Next()
}
