package handler

import (
	"films/internal/models"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"net/http"
)

type AuthService interface {
	SignUp(user *models.User) error
	SignIn(user *models.User) (string, error)
	ParseToken(accessToken string) (*models.User, error)

	SendVerificationEmail(user *models.User) error
	VerifyEmail(email, hash string) error
	ResetPassword(user *models.User) error
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(a AuthService) *AuthHandler {
	return &AuthHandler{service: a}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	inp := new(models.User)

	err := c.BindJSON(inp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// validating
	err = validation.ValidateStruct(inp,
		validation.Field(&inp.FullName, validation.Required),
		validation.Field(&inp.Email, validation.Required, is.Email),
		validation.Field(&inp.Password, validation.Required),
		validation.Field(&inp.Type, validation.Required, validation.In("user", "company")),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// signing up
	err = h.service.SignUp(inp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// sending verification email
	err = h.service.SendVerificationEmail(inp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	inp := new(models.User)

	err := c.BindJSON(inp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// validation
	err = validation.ValidateStruct(inp,
		validation.Field(&inp.Email, validation.Required, is.Email),
		validation.Field(&inp.Password, validation.Required),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// returning token
	token, err := h.service.SignIn(inp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, signInResponse{Token: token})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	email := c.Query("email")
	hash := c.Query("hash")

	err := h.service.VerifyEmail(email, hash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(200, "email verified")
}

func (h *AuthHandler) ResetPassword(c *gin.Context){
	user := new(models.User)
	err := c.BindJSON(&user)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.ResetPassword(user)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}