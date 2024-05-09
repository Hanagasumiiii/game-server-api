package user

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

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

func (s *Service) CreateUser(username, email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err = s.db.Exec(query, username, email, hashedPassword)
	return err
}

func (s *Service) AuthenticateUser(email, password string) (bool, error) {
	var hashedPassword string
	query := `SELECT password FROM users WHERE username = $1`
	err := s.db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}
	return CheckPassword(hashedPassword, password), nil
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
