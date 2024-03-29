package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"aldiyarsagidolla/snippetbox/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = u.DB.Exec(stmt, name, email, string(hashedPw))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message,
				"users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil

}

func (u *UserModel) Authenticate(email, password string) (int, error) {

	var id int
	var hashedPw []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE`
	row := u.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPw, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (u *UserModel) Get(id int) (*models.User, error) {
	usr := &models.User{}

	stmt := `SELECT id, name, email, created, active FROM users WHERE id = ?`
	err := u.DB.QueryRow(stmt, id).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Created, &usr.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return usr, nil
}
