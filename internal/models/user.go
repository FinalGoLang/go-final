package models

import "time"

type User struct {
	UserID      int       `db:"user_id" json:"user_id,omitempty"`
	FullName    string    `db:"full_name" json:"full_name,omitempty"`
	Phone       string    `db:"phone" json:"phone,omitempty"`
	Email       string    `db:"email" json:"email,omitempty"`
	OldPassword string    `json:"old_password,omitempty"`
	Password    string    `db:"password" json:"password,omitempty"`
	Verified    bool      `db:"verified" json:"verified,omitempty"`
	Hash        string    `db:"hash" json:"hash,omitempty"`
	Type        string    `db:"type" json:"type"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}
