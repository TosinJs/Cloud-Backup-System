package mySqlRepo

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"tosinjs/cloud-backup/internal/entity/errorEntity"
	"tosinjs/cloud-backup/internal/entity/userEntity"
	"tosinjs/cloud-backup/internal/repository/userRepo"
)

type mySql struct {
	conn *sql.DB
}

func New(conn *sql.DB) userRepo.UserRepository {
	return mySql{
		conn: conn,
	}
}

func (m mySql) CreateUser(req userEntity.UserSignUpReq) *errorEntity.ServiceError {
	stmt := fmt.Sprintf(`
		INSERT INTO Users(user_id, username, password, email) 
		VALUES('%v', '%v', '%v', '%v')
		`, req.UserId, req.Username, req.Password, req.Email)

	_, err := m.conn.Exec(stmt)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 {
				return errorEntity.ConflictError("Duplicate Email or Username", err)
			}
		}
		return errorEntity.InternalServerError(err)
	}
	return nil
}

func (m mySql) LoginUser(req userEntity.UserLoginReq) (string, *errorEntity.ServiceError) {
	stmt := fmt.Sprintf(`
		SELECT password FROM Users WHERE username = '%v'
		`, req.Username)

	var password string
	err := m.conn.QueryRow(stmt).Scan(&password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", errorEntity.NotFoundError("Invalid Login Credentials", err)
		}
		return "", errorEntity.InternalServerError(err)
	}
	return password, nil
}
