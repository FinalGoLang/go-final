package service

import (
	"errors"
	"films/internal/models"
)

type FilmStorage interface {
	GetByID(id int) (*models.Films, error)
	GetAll(id int) (*[]models.Films, error)
	Insert(id int, film *models.Films) error
	Update(film *models.Films, id int) error
	Delete(id, filmID int) error

	GetAllFilms() (*[]models.Films, error)
	GetBestFilms() (*[]models.Films, error)
}

type Film struct {
	repo FilmStorage
}

func NewFilmService(storage FilmStorage) *Film {
	return &Film{repo: storage}
}

func (f *Film) GetBestFilms() (*[]models.Films, error){

	films, err := f.repo.GetBestFilms()
	if err != nil{
		return nil, err
	}

	return films, nil
}

func (f *Film) GetAllFilms() (*[]models.Films, error){

	films, err := f.repo.GetAllFilms()
	if err != nil{
		return nil, err
	}

	return films, nil
}

func (f *Film) Create(user *models.User, film *models.Films) error{
	if user.Type != "company"{
		return errors.New("only users with company type can add films")
	}
	err := f.repo.Insert(user.UserID, film)
	if err != nil{
		return err
	}

	return nil
}

func (f *Film) GetByID(filmID int) (*models.Films, error){
	film, err := f.repo.GetByID(filmID)
	if err != nil{
		return nil, err
	}

	return film ,err
}

func (f *Film) GetAll(user *models.User) (*[]models.Films, error){
	if user.Type != "company"{
		return nil, errors.New("only users with company type can get all his/her films")
	}
	films, err := f.repo.GetAll(user.UserID)
	if err != nil{
		return nil, err
	}

	return films ,err
}

func (f *Film) Update(user *models.User, film *models.Films) error{
	if user.Type != "company"{
		return errors.New("only users with company type can update his/her films")
	}

	err := f.repo.Update(film, user.UserID)
	if err != nil{
		return err
	}

	return nil
}

func (f *Film) Delete(user *models.User, filmID int) error{
	if user.Type != "company"{
		return errors.New("only users with company type can delete his/her films")
	}
	err := f.repo.Delete(user.UserID, filmID)
	if err != nil{
		return err
	}

	return nil
}