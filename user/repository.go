package user

import (
	"testZaShtat/models"
)

type Repository interface {
	GetAllUsers() ([]*User, error)
	GetByID(id rune) (*User, error)
	Store(u *User) error
	Update(u *User) error
	Delete(id rune) error
}

type UserRepository struct {
	DB *models.DB
}

func (r *UserRepository) GetAllUsers() ([]*User, error) {
	rows, err := r.DB.Query("SELECT id, first_name, middle_name, last_name, username FROM \"user\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		u := new(User)
		err := rows.Scan(&u.ID, &u.FirstName, &u.MiddleName, &u.LastName, &u.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetByID(id rune) (*User, error) {
	stmt, err := r.DB.Prepare("SELECT id, first_name, middle_name, last_name, username FROM \"user\" WHERE id=$1")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(User)
	err = rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Store(u *User) error {
	stmt, err := r.DB.Prepare("INSERT INTO \"user\" (first_name, middle_name, last_name, username, password_hash)  VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	_, err = stmt.Query(u.FirstName, u.MiddleName, u.LastName, u.Username, u.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(u *User) error {
	stmt, err := r.DB.Prepare("UPDATE \"user\" SET first_name=$1, middle_name=$2, last_name=$3, username=$4) WHERE id=$6")
	if err != nil {
		return err
	}
	_, err = stmt.Query(u.FirstName, u.MiddleName, u.LastName, u.Username, u.ID)
	if err != nil {
		return err
	}
	return nil
}
