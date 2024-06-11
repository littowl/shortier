package db

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

type Link struct {
	Id   int
	Link string `json:"link"`
	Hash string
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}

func DbStart(baseUrl string) *pgxpool.Pool {
	fmt.Print(baseUrl)
	dbpool, err := pgxpool.New(context.Background(), baseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v", err)
		os.Exit(1)
	}
	return dbpool
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateHash() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (db DB) CreateHash(link Link) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	hash := generateHash()
	_, err = conn.Exec(context.Background(), "INSERT INTO links(link, hash) VALUES ($1, $2)", link.Link, hash)
	// err = conn.QueryRow(context.Background(),
	// 	"INSERT INTO links(link, hash) VALUES ($1, $2) RETURNING id", link.Link, hash).Scan(&link.Id)

	if err != nil {
		return fmt.Errorf("unable to INSERT: %v", err)
	}

	return nil
}

func (db DB) GetById(id int) (*string, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	row, err := conn.Query(context.Background(), "SELECT * FROM links WHERE id = $1", id)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from database: %v", err)
	}
	defer row.Close()

	var l Link
	var shortLink string

	for row.Next() {
		err = row.Scan(&l.Id, &l.Link, &l.Hash)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		shortLink = strings.Split(l.Link, "/")[0] + "/" + l.Hash
	}

	return &shortLink, nil
}
