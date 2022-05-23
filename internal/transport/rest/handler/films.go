package handler

import (
	"films/internal/models"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type FilmService interface {
	Create(user *models.User, film *models.Films) error
	GetByID(filmID int) (*models.Films, error)
	GetAll(user *models.User) (*[]models.Films, error)
	Update(user *models.User, film *models.Films) error
	Delete(user *models.User, filmID int) error

	GetAllFilms()(*[]models.Films, error)
	GetBestFilms()(*[]models.Films, error)
}

type FilmHandler struct {
	service FilmService
}

func NewFilmHandler(u FilmService) *FilmHandler {
	return &FilmHandler{service: u}
}

func (h *FilmHandler) GetBestFilms(c *gin.Context) {
	bestFilms, err := h.service.GetBestFilms()
	if err != nil{
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, bestFilms)
}

func (h *FilmHandler) GetAllFilms(c *gin.Context) {
	allFilms, err := h.service.GetAllFilms()
	if err != nil{
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return

	}

	c.JSON(http.StatusOK, allFilms)
}

func (h *FilmHandler) Create(c *gin.Context){
	user := c.MustGet(CtxUserKey).(*models.User)

	film := new(models.Films)
	err := c.BindJSON(&film)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = validation.ValidateStruct(film,
		validation.Field(&film.Name,validation.Required),
		validation.Field(&film.Price, validation.Required),
		validation.Field(&film.Rating, validation.Required, validation.In(1,2,3,4,5,0)),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Create(user, film)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *FilmHandler) GetAllCompanyFilms(c *gin.Context) {
	user := c.MustGet(CtxUserKey).(*models.User)

	films, err := h.service.GetAll(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, films)
}

func (h *FilmHandler) GetFilmByID(c *gin.Context) {
	inp := new(models.Films)
	err := c.BindJSON(&inp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = validation.ValidateStruct(inp, validation.Field(&inp.FilmID,validation.Required))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	films, err := h.service.GetByID(inp.FilmID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, films)
}

func (h *FilmHandler) Update(c *gin.Context) {
	updFilm := new(models.Films)
	err := c.BindJSON(&updFilm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = validation.ValidateStruct(updFilm,
		validation.Field(&updFilm.FilmID,validation.Required),
		validation.Field(&updFilm.Name,validation.Required),
		validation.Field(&updFilm.Price, validation.Required),
		validation.Field(&updFilm.Rating, validation.Required),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet(CtxUserKey).(*models.User)

	err = h.service.Update(user, updFilm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *FilmHandler) Delete(c *gin.Context) {
	user := c.MustGet(CtxUserKey).(*models.User)

	film := new(models.Films)
	err := c.BindJSON(&film)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// validating
	err = validation.ValidateStruct(film,
		validation.Field(&film.FilmID,validation.Required),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// deleting user's film
	err = h.service.Delete(user, film.FilmID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
