package user

import "database/sql"

type Service struct {
	db *sql.DB
}

type User struct {
	ID       int
	Username string
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetUserByID(id int) (*User, error) {
	var user User
	query := `SELECT id, username FROM users WHERE id = $1`
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUserByUsername(username string) (*User, error) {
	var user User
	query := `SELECT id, username FROM users WHERE username = $1`
	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) UpdateUser(user *User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	query := `UPDATE users SET username = $1 WHERE id = $2`
	_, err = tx.Exec(query, user.Username, user.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
