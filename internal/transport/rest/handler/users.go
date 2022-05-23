package handler

import (
	"films/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"net/http"
)

type UserService interface {
	Create(user *models.User) (*models.User, error)
	GetByID(userID int) (*models.User, error)
	GetAll() (*[]models.User, error)
	Update(user *models.User) error
	Delete(userID int) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(u UserService) *UserHandler {
	return &UserHandler{service: u}
}


func (h *UserHandler) GetUserInfo(c *gin.Context) {
	user := c.MustGet(CtxUserKey).(*models.User)

	user, err := h.service.GetByID(user.UserID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	updUser := new(models.User)
	err := c.BindJSON(&updUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = validation.ValidateStruct(updUser,
		validation.Field(&updUser.Email,validation.Required, is.Email),
		validation.Field(&updUser.Phone, validation.Required),
		validation.Field(&updUser.FullName, validation.Required),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	oldUser := c.MustGet(CtxUserKey).(*models.User)
	updUser.UserID = oldUser.UserID

	err = h.service.Update(updUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *UserHandler) Delete(c *gin.Context) {
	user := c.MustGet(CtxUserKey).(*models.User)

	err := h.service.Delete(user.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
