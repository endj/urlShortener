package main

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
var getSetStatement *sql.Stmt
var findStatement *sql.Stmt

func prepareStatements() error {
	insert := "CALL GetOrSet(?, ?, ?)"
	find := "SELECT url FROM urls WHERE shortUrl = ?"

	var err error
	getSetStatement, err = db.Prepare(insert)
	if err != nil {
		return err
	}

	findStatement, err = db.Prepare(find)
	if err != nil {
		return err
	}

	return err
}

func initTables(db *sql.DB) error {
	urlTable := `CREATE TABLE IF NOT EXISTS urls(
        shortUrl VARBINARY(7),
        hash VARBINARY(32),
		url VARBINARY(2048),
        INDEX idx_hash (hash),
		INDEX url_hash (shortUrl)
    )`
	_, err := db.Exec(urlTable)
	return err
}

func insertProcExists() (bool, error) {
	procedureCheckSQL := `
	SELECT COUNT(*) FROM information_schema.ROUTINES
	WHERE ROUTINE_TYPE = 'PROCEDURE' AND ROUTINE_SCHEMA = 'shorturl' AND ROUTINE_NAME = 'GetOrSet'
	`
	var count int
	err := db.QueryRow(procedureCheckSQL).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createInsertProcedure() error {
	exists, err := insertProcExists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	var insert = `
	CREATE PROCEDURE GetOrSet(
		IN in_shortUrl VARBINARY(7),
		IN in_hash VARBINARY(32),
		IN in_url VARBINARY(2048)
	)
	BEGIN
		DECLARE out_url VARBINARY(7);
	
		INSERT INTO urls (shortUrl, hash, url)
		VALUES (in_shortUrl, in_hash, in_url)
		ON DUPLICATE KEY UPDATE url = IF(hash = in_hash, VALUES(url), url);
	
		SELECT shortUrl INTO out_url FROM urls WHERE hash = in_hash LIMIT 1;
	
		SELECT out_url;
	END;	
	`
	_, err = db.Exec(insert)
	return err
}

func SetupDBConnection() error {
	var err error
	db, err = sql.Open("mysql", "shorturl:shorturl@tcp(localhost:3306)/shorturl")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	err = initTables(db)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	err = createInsertProcedure()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	err = prepareStatements()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
