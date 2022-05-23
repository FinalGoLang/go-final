package psql

import (
	"films/internal/models"
	"github.com/jmoiron/sqlx"
)

type Film struct {
	db *sqlx.DB
}

func NewFilm(db *sqlx.DB) *Film {
	return &Film{db}
}
func (f *Film) GetBestFilms() (*[]models.Films, error) {
	var films []models.Films
	query := `SELECT * FROM films WHERE rating>4;`

	err := f.db.Select(&films, query)
	if err != nil {
		return nil, err
	}

	return &films, nil
}

func (f *Film) GetAllFilms() (*[]models.Films, error) {
	var films []models.Films
	query := `SELECT * FROM films;`

	err := f.db.Select(&films, query)
	if err != nil {
		return nil, err
	}

	return &films, nil
}

func (f *Film) GetByID(id int) (*models.Films, error) {
	var film models.Films
	query := `SELECT * FROM films WHERE film_id=$1`

	err := f.db.Get(&film, query, id)
	if err != nil {
		return nil, err
	}

	return &film, nil
}

// GetAll gets all company films
func (f *Film) GetAll(id int) (*[]models.Films, error) {
	var films []models.Films
	query := `SELECT * FROM films WHERE user_id=$1;`

	err := f.db.Select(&films, query, id)
	if err != nil {
		return nil, err
	}

	return &films, nil
}
func (f *Film) Insert(id int, film *models.Films) error {
	query := `INSERT INTO films (name, price, rating, user_id) VALUES (:name,:price,:rating, :userID);`
	queryArgs := map[string]interface{}{"name": film.Name, "price": film.Price, "rating": film.Rating, "userID": id}

	_, err := f.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}
func (f *Film) Update(film *models.Films, id int) error {
	query := `
		UPDATE films 
		SET name=:name,
			price=:price,
		    rating=:rating,
		    updated_at=current_timestamp
		WHERE user_id=:userID and film_id=:filmID;
		`
	queryArgs := map[string]interface{}{"userID": id, "name": film.Name, "price": film.Price, "rating": film.Rating, "filmID": film.FilmID}

	_, err := f.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}
func (f *Film) Delete(id, filmID int) error {
	query := `DELETE FROM films WHERE user_id=:id and film_id=:filmID`
	queryArgs := map[string]interface{}{"id": id, "filmID": filmID}

	_, err := f.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}
