package handler

import (
	"films/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SubscribeService interface {
	Subscribe(user *models.User, film *models.Films) error
	GetMySubscriptions(user *models.User) (*[]models.MyFilms, error)
	GetExpiredSubscribes(user *models.User, now time.Time) (*[]models.MyFilms, error)
}

type SubscribeHandler struct {
	service SubscribeService
}

func NewSubscribeHandler(s SubscribeService) *SubscribeHandler{
	return &SubscribeHandler{s}
}

func (s *SubscribeHandler) Subscribe(c *gin.Context){
	film := new(models.Films)
	err := c.BindJSON(&film)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}


	user := c.MustGet(CtxUserKey).(*models.User)
	err = s.service.Subscribe(user,film)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (s *SubscribeHandler) MySubscriptions(c *gin.Context){
	user := c.MustGet(CtxUserKey).(*models.User)

	films, err := s.service.GetMySubscriptions(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, films)
}

func (s *SubscribeHandler) GetExpiredSubscribes(c *gin.Context){
	user := c.MustGet(CtxUserKey).(*models.User)
	now := time.Now()

	films, err := s.service.GetExpiredSubscribes(user,now)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, films)
}