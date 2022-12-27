package mySqlRepo

import (
	"database/sql"
	"fmt"
	"tosinjs/cloud-backup/internal/entity/errorEntity"
	// "github.com/go-sql-driver/mysql"
	"tosinjs/cloud-backup/internal/repository/fileRepo"
)

type mySql struct {
	conn *sql.DB
}

func New(conn *sql.DB) fileRepo.FileRepository {
	return mySql{
		conn: conn,
	}
}

func (m mySql) UploadFile(username, filename string) *errorEntity.ServiceError {
	stmt := fmt.Sprintf(`
		INSERT INTO Media(username, filepath) 
		VALUES('%v', '%v')
		`, username, filename)

	_, err := m.conn.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return errorEntity.InternalServerError(err)
	}
	return nil
}
func (m mySql) DeleteFile(filename string) *errorEntity.ServiceError {
	stmt := fmt.Sprintf(`
		DELETE FROM Media
		WHERE filepath = '%v'
		`, filename)

	_, err := m.conn.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return errorEntity.InternalServerError(err)
	}
	return nil
}

func (m mySql) FlagFile(filename string) (int, *errorEntity.ServiceError) {
	stmt := fmt.Sprintf(`
		UPDATE Media SET flag_count = flag_count + 1
		WHERE filepath = '%v'
	`, filename)

	_, err := m.conn.Exec(stmt)
	if err != nil {
		return 0, errorEntity.InternalServerError(err)
	}

	stmt = fmt.Sprintf(`
		SELECT flag_count FROM Media
		WHERE filepath = '%v'
	`, filename)

	var flagCount int
	err = m.conn.QueryRow(stmt).Scan(&flagCount)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, errorEntity.NotFoundError("File Not Found", err)
		}
		return 0, errorEntity.InternalServerError(err)
	}

	return flagCount, nil
}
