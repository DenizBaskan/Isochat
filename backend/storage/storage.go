package storage

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	conn *sql.Conn
}

type CacheDB struct {
	client *redis.Client
}

var (
	Cache CacheDB
	DB    Database
)

func Setup() {
	// create a new redis client
	Cache = CacheDB{client: redis.NewClient(&redis.Options{})}

	// connect to the mysql database (dsn from .env)
	d, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}

	conn, err := d.Conn(context.TODO())
	if err != nil {
		panic(err)
	}

	DB = Database{conn: conn}
}

func (c *CacheDB) Set(key string, value any, expiration time.Duration) error {
	return Cache.client.Set(key, value, expiration).Err()
}

// returns empty string if key does not exist
func (c *CacheDB) Get(key string) (string, error) {
	v, err := Cache.client.Exists(key).Result()
	if err != nil {
		return "", err
	}

	if v == 1 {
		v, err := Cache.client.Get(key).Result()
		if err != nil {
			return "", err
		}

		return v, nil
	}

	return "", nil
}

// does a record exist?
func (db Database) Exists(fmt string, args ...any) (bool, error) {
	res := db.conn.QueryRowContext(context.TODO(), "SELECT 1 FROM " + fmt, args...)
	if err := res.Err(); err != nil {
		return false, err
	}
	
	var tmp int64

	if err := res.Scan(&tmp); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// insert a record into a table
func (db Database) Insert(fmt string, args ...any) error {
	_, err := db.conn.ExecContext(context.TODO(), "INSERT INTO " + fmt, args...)
	return err
}

// get ONE row and store result in scanArgs
func (db Database) GetRow(scanArgs []any, fmt string, args ...any) error {
	res := db.conn.QueryRowContext(context.TODO(), "SELECT " + fmt, args...)
	if err := res.Err(); err != nil {
		return err
	}

	return res.Scan(scanArgs...)
}

// get rows
func (db Database) GetRows(scanArgs []any, fmt string, args ...any) ([][]any, error) {
	var faces [][]any

	rows, err := db.conn.QueryContext(context.TODO(), "SELECT " + fmt, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		cpy := scanArgs
		rows.Scan(cpy...)
		faces = append(faces, cpy)
	}

	return faces, nil
}

// delete row
func (db Database) Delete(fmt string, args ...any) error {
	_, err := db.conn.ExecContext(context.TODO(), "DELETE FROM " + fmt, args...)
	return err
}

// update a row
func (db Database) Update(fmt string, args ...any) error {
	_, err := db.conn.ExecContext(context.TODO(), "UPDATE " + fmt, args...)
	return err
}
