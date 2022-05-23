package service

import (
	"errors"
	"films/internal/models"
	"time"
)

type SubscribeStorage interface {
	Insert(userID, filmID int,expires time.Time) error
	GetSubscription(userID int) (*[]models.MyFilms, error)
	GetExpiredSubscription(userID int, now time.Time)(*[]models.MyFilms, error)
}

type Subscribe struct {
	repo SubscribeStorage
}

func NewSubscribeService(repo SubscribeStorage) *Subscribe {
	return &Subscribe{repo}
}

func (s *Subscribe) Subscribe(user *models.User, film *models.Films) error {
	if user.Type != "user" {
		return errors.New("only users with user type can add subcriptions")
	}
	err := s.repo.Insert(user.UserID, film.FilmID, time.Now().Add(time.Hour * 360))
	if err != nil {
		return err
	}
	return nil
}
func (s *Subscribe) GetMySubscriptions(user *models.User) (*[]models.MyFilms, error) {
	if user.Type != "user" {
		return nil, errors.New("only users with user type can get their subcriptions")
	}

	films, err := s.repo.GetSubscription(user.UserID)
	if err != nil {
		return nil, err
	}

	return films, nil
}

func (s *Subscribe) GetExpiredSubscribes(user *models.User, now time.Time) (*[]models.MyFilms, error){
	if user.Type != "user" {
		return nil, errors.New("only users with user type can get their subcriptions")
	}

	films, err := s.repo.GetExpiredSubscription(user.UserID, now)
	if err != nil {
		return nil, err
	}

	return films, nil
}