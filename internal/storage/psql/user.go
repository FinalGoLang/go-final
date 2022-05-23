package psql

import (
	"films/internal/models"
	"github.com/jmoiron/sqlx"
)

type User struct {
	db *sqlx.DB
}

func NewAccount(db *sqlx.DB) *User {
	return &User{db}
}

func (r *User) Insert(user *models.User) error {
	query := `INSERT INTO users (full_name, phone, email, password, hash,type) VALUES (:fullName, :phone, :email, :password, :hash, :type);`
	queryArgs := map[string]interface{}{"fullName": user.FullName, "phone": user.Phone, "email": user.Email, "password": user.Password, "hash": user.Hash,"type":user.Type}

	_, err := r.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}

func (r *User) GetAll() (*[]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users;`

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *User) GetByCredentials(email, password string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email=$1 and password=$2`

	err := r.db.Get(&user, query, email, password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *User) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE user_id=$1`

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *User) Update(user *models.User, id int) error {
	query := `
		UPDATE users 
		SET full_name=:fullName, 
		    phone=:phone, 
		    email=:email,
		    updated_at=current_timestamp
		WHERE user_id=:userID;
		`
	queryArgs := map[string]interface{}{"userID": id, "fullName": user.FullName, "phone": user.Phone, "email": user.Email}

	_, err := r.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}

func (r *User) Delete(id int) error {
	query := `DELETE FROM users WHERE user_id=:id`
	queryArgs := map[string]interface{}{"id": id}

	_, err := r.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}

func (r *User) GetUserHash(email string) (string, error) {
	var hash string
	query := `SELECT hash FROM users WHERE email=$1`

	err := r.db.Get(&hash, query, email)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (r *User) VerifyUser(email string) error {
	query := `
		UPDATE users 
		SET verified=true
		WHERE email=:email;
		`
	queryArgs := map[string]interface{}{"email": email}

	_, err := r.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}

func (r *User) GetPasswordHash(password string) (string, error) {
	var hash string
	query := `SELECT hash FROM users WHERE password=$1`

	err := r.db.Get(&hash, query, password)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (r *User) UpdateUserPassword(user *models.User) error {
	query := `
		UPDATE users 
		SET password=:password
		WHERE email=:email;
		`
	queryArgs := map[string]interface{}{"email": user.Email, "password":user.Password}

	_, err := r.db.NamedExec(query, queryArgs)
	if err != nil {
		return err
	}

	return nil
}