package psql

import (
	"films/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Subscribe struct {
	db *sqlx.DB
}

func NewSubscribe(db *sqlx.DB) *Subscribe{
	return &Subscribe{db}
}

func (s *Subscribe) Insert(userID, filmID int, expire time.Time) error{
	query := `INSERT INTO subscriptions (user_id, film_id, expires) VALUES (:userID, :filmID, :expire);`
	queryArgs := map[string]interface{}{"userID": userID, "filmID":filmID, "expire":expire}

	_, err := s.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}
func (s *Subscribe) GetSubscription(userID int) (*[]models.MyFilms, error) {
	var users []models.MyFilms
	query := `
	SELECT f.film_id, f.name, f.price,s.expires
	FROM users u, films f, subscriptions s
	WHERE u.user_id = s.user_id and f.film_id = s.film_id and u.user_id=$1 and s.expires>current_timestamp;`

	err := s.db.Select(&users, query, userID)
	if err != nil {
		return nil, err
	}

	return &users, nil
}
func (s *Subscribe) GetExpiredSubscription(userID int, now time.Time)(*[]models.MyFilms, error) {
	var users []models.MyFilms
	query := `
	SELECT f.film_id, f.name, f.price,s.expires
	FROM users u, films f, subscriptions s
	WHERE u.user_id = s.user_id and f.film_id = s.film_id and u.user_id=$1 and s.expires<$2;`

	err := s.db.Select(&users, query, userID, now)
	if err != nil {
		return nil, err
	}

	return &users, nil
}