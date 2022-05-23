package service

import (
	"films/internal/models"
)

type UserStorage interface {
	GetByCredentials(email, password string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetAll() (*[]models.User, error)
	Insert(user *models.User) error
	Update(user *models.User, id int) error
	Delete(id int) error

	GetUserHash(email string) (string, error)
	VerifyUser(email string) error
	GetPasswordHash(password string) (string, error)
	UpdateUserPassword(user *models.User) error
}

type User struct {
	repo UserStorage
}

func NewUserService(storage UserStorage) *User {
	return &User{repo: storage}
}

func (s *User) Create(user *models.User) (*models.User, error) {
	err := s.repo.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *User) GetByID(userID int) (*models.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *User) GetAll() (*[]models.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *User) Update(user *models.User) error {

	err := s.repo.Update(user, user.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *User) Delete(userID int) error {
	err := s.repo.Delete(userID)
	if err != nil {
		return err
	}

	return nil
}
