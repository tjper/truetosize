// Package shoesdb provides a library for interacting with the shoes Postgres database.
package shoesdb

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// Error codes returned by failures within library.
var (
	ErrZeroShoes             = errors.New("Zero shoes passed to InsertShoes.")
	ErrZeroTrueToSizes       = errors.New("Zero truetosizes passed to InsertTrueToSizes.")
	ErrIdentifierInvalidType = errors.New("SelectTrueToSize identifier must be an int or a string data type.")
)

var (
	logErr  *log.Logger
	logInfo *log.Logger
)

// init initializes the loggers logErr and logInfo.
func init() {
	file := "/home/james/go/log/shoesdb.txt"
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logErr = log.New(f, "Error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	logInfo = log.New(f, "Error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
}

// Conn establishes a connection with the shoes db, if successful,
// returns db connection.
func Conn() *sql.DB {
	connStr := "user=shoes dbname=shoes sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logErr.Fatalln(err)
	}
	return db
}

type Execer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

type Queryer interface {
	Query(string, ...interface{}) (Rower, error)
}

type Rower interface {
	Close() error
	Err() error
	Next() bool
	Scan(...interface{}) error
}

// InsertShoes inserts a new shoe row into a Execer, if successful,
// returns number of rows inserted.
func InsertShoes(db Execer, shoes ...string) (int, error) {
	if len(shoes) == 0 {
		logErr.Println(ErrZeroTrueToSizes)
		return 0, ErrZeroShoes
	}

	// write query string
	query := `INSERT INTO shoes (name) VALUES`
	values := make([]string, len(shoes))
	for i, _ := range values {
		values[i] = "(?)"
	}
	strings.Join(values, ",")
	query += ";"

	// convert []int to []interface{}
	iShoes := make([]interface{}, len(shoes))
	for i, v := range shoes {
		iShoes[i] = v
	}

	res, err := db.Exec(query, iShoes...)
	if err != nil {
		logErr.Println(err)
		return 0, err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		logErr.Println(err)
		return 0, err
	}
	return int(aff), nil
}

// InsertTrueToSize inserts a new truetosize row into a Execer, if successful,
// returns the number of rows inserted.
func InsertTrueToSizes(db Execer, trueToSizes ...int) (int, error) {
	if len(trueToSizes) == 0 {
		logErr.Println(ErrZeroTrueToSizes)
		return 0, ErrZeroTrueToSizes
	}

	// write query string
	query := `INSERT INTO truetosize (truetosize) VALUES`
	values := make([]string, len(trueToSizes))
	for i, _ := range values {
		values[i] = "(?)"
	}
	strings.Join(values, ",")
	query += ";"

	// convert []int to []interface{}
	iTrueToSizes := make([]interface{}, len(trueToSizes))
	for i, v := range trueToSizes {
		iTrueToSizes[i] = v
	}

	res, err := db.Exec(query, iTrueToSizes...)
	if err != nil {
		logErr.Println(err)
		return 0, err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		logErr.Println(err)
		return 0, err
	}
	return int(aff), nil
}

// SelectTrueToSizeByShoesId retrieves each truetosize rating from Queryer by identifiear,
// identifier may be a shoes' id or a shoes' name, if successful,
// returns the set of truetosize ratings.
func SelectTrueToSize(db Queryer, identifier interface{}) ([]int, error) {
	query := `SELECT t.truetosize
		  FROM truetosize t`

	if _, ok := identifier.(int); ok {
		query += `WHERE shoes_id = ?`
	} else if _, ok := identifier.(string); ok {
		query += `INNER JOIN shoes s
			  ON (t.shoes_id = s.id)
			  WHERE s.name = ?;`
	} else {
		logErr.Println(ErrIdentifierInvalidType)
		return nil, ErrIdentifierInvalidType
	}

	rows, err := db.Query(query, identifier)
	if err != nil {
		logErr.Println(err)
		return nil, err
	}
	defer rows.Close()

	ttsSet := make([]int, 0)
	for rows.Next() {
		var truetosize int
		if err := rows.Scan(&truetosize); err != nil {
			logErr.Println(err)
			return nil, err
		}
		ttsSet = append(ttsSet, truetosize)
	}
	if err := rows.Err(); err != nil {
		logErr.Println(err)
		return nil, err
	}
	return ttsSet, nil
}
