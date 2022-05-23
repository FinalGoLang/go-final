package handler

import (
	"github.com/gin-gonic/gin"
)

type Deps struct {
	UserHandler      UserHandler
	AuthHandler      AuthHandler
	FilmHandler      FilmHandler
	SubscribeHandler SubscribeHandler
}
type Handler struct {
	UserHandler      *UserHandler
	AuthHandler      *AuthHandler
	FilmHandler      *FilmHandler
	SubscribeHandler *SubscribeHandler
}

func New(dep Deps) *Handler {
	return &Handler{
		UserHandler:      &dep.UserHandler,
		AuthHandler:      &dep.AuthHandler,
		FilmHandler:      &dep.FilmHandler,
		SubscribeHandler: &dep.SubscribeHandler,
	}
}

func (h *Handler) InitRoutes(service AuthService) *gin.Engine {
	app := gin.Default()

	app.POST("/sign-up", h.AuthHandler.SignUp)
	app.POST("/sign-in", h.AuthHandler.SignIn)
	app.POST("/reset-pwd", h.AuthHandler.ResetPassword)
	app.GET("/verify", h.AuthHandler.VerifyEmail)

	app.GET("/films", h.FilmHandler.GetAllFilms)
	app.GET("/best", h.FilmHandler.GetBestFilms)
	app.GET("/get", h.FilmHandler.GetFilmByID)

	middleware := NewAuthMiddleware(service)
	api := app.Group("/api", middleware)

	users := api.Group("/users")
	users.GET("", h.UserHandler.GetUserInfo)
	users.PUT("", h.UserHandler.Update)
	users.DELETE("", h.UserHandler.Delete)

	films := api.Group("/films")

	films.POST("", h.FilmHandler.Create)
	films.GET("", h.FilmHandler.GetAllCompanyFilms)
	films.PUT("", h.FilmHandler.Update)
	films.DELETE("", h.FilmHandler.Delete)



	subscribe := users.Group("/subscribe")
	subscribe.POST("", h.SubscribeHandler.Subscribe)
	subscribe.GET("", h.SubscribeHandler.MySubscriptions)
	subscribe.GET("/expired", h.SubscribeHandler.GetExpiredSubscribes)

	return app
}
