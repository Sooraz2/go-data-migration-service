package services

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-data-migration/models"
)

func ConnectDB(config models.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	return sql.Open("mysql", dsn)
}

func TableExists(db *sql.DB, tableName string) (bool, error) {
	query := "SHOW TABLES LIKE ?"
	var table string
	err := db.QueryRow(query, tableName).Scan(&table)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func CreateTable(destDB *sql.DB, sourceDB *sql.DB, tableName string) error {
	// Get CREATE TABLE statement from source
	var tableCreate string
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	err := sourceDB.QueryRow(query).Scan(&tableName, &tableCreate)
	if err != nil {
		return err
	}

	// Execute CREATE TABLE in destination
	_, err = destDB.Exec(tableCreate)
	return err
}
