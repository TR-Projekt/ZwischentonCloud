package token

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ADMIN       int = 42
	CREATOR     int = 1
	COORDINATOR int = 2
)

type User struct {
	ID           int       `json:"user_id" sql:"user_id"`
	Email        string    `json:"user_email" sql:"user_email"`
	PasswordHash string    `json:"user_password" sql:"user_password"`
	CreateDate   time.Time `json:"user_createdat" sql:"user_createdat"`
	UpdateDate   time.Time `json:"user_updatedat" sql:"user_updatedat"`
	Role         int       `json:"user_role" sql:"user_role"`
}

func UserScan(rs *sql.Rows) (User, error) {
	var u User
	return u, rs.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreateDate, &u.UpdateDate, &u.Role)
}

type UserClaims struct {
	UserID           string
	UserRole         int
	UserZwischentons []int
	UserSituations   []int
	jwt.RegisteredClaims
}
