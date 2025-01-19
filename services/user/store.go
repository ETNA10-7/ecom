package user

import (
	"database/sql"
	"fmt"

	"github.com/ETNA10-7/ecom/types"
	"github.com/lib/pq"
	//"golang.org/x/tools/go/analysis/passes/defers"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//Its a pointer with zeroed initialization
	//Also can be written as &types.User
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		// if u.ID == 0 {
		// 	return nil, fmt.Errorf("User not found")
		// }
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("User not found")
	}
	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("User Not Found")
	}
	return u, nil

}
func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES ($1,$2,$3,$4)", user.FirstName, user.LastName, user.Email, user.Password)
	// if err != nil {
	// 	return err
	// }
	if err != nil {
		// Check if the error is a unique constraint violation
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// PostgreSQL error code 23505 is for unique_violation
			return fmt.Errorf("email already exists")
		}
		return err // Handle other errors
	}
	return nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	// The below is a zeroed pointer means does not have any values
	// It can also be written as user:=&types.User
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
