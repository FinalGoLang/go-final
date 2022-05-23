package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const CtxUserKey = "user"

type AuthMiddleware struct {
	authService AuthService
}

func NewAuthMiddleware(authService AuthService) gin.HandlerFunc {
	return (&AuthMiddleware{
		authService: authService,
	}).Handle
}
func (m *AuthMiddleware) Handle(c *gin.Context){
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.authService.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Set(CtxUserKey, user)
}
