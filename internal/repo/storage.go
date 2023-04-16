package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Albitko/loyalty-program/internal/entities"
)

const (
	uniqueViolationErr = "23505"
)
const schema = `
 	CREATE TABLE IF NOT EXISTS users (
		id text primary key,
		login text not null unique,
		password text not null
	);
	CREATE TABLE IF NOT EXISTS orders (
	  	"order" text primary key unique,
	  	user_id text not null references users(id),
	    status text not null,
	    accrual float not null,
		"current" float not null default 0,
        withdrawn float not null  default 0,
	    uploaded_at timestamp
	);
	CREATE TABLE IF NOT EXISTS withdrawals (
	    user_id text not null references users(id),
		"sum" float not null,
		processed_at timestamp
	);
 	`

type repository struct {
	db  *sql.DB
	ctx context.Context
}

func (r *repository) UpdateOrder(ctx context.Context, order entities.Order) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetUserBalance(ctx context.Context, user string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetUserWithdrawn(ctx context.Context, user string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Withdraw(ctx context.Context, amount string) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetUserForOrder(ctx context.Context, order string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) CreateOrder(ctx context.Context, order string) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetOrdersForUser(ctx context.Context, user string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Register(ctx context.Context, id, login, hashedPassword string) error {
	var pgErr *pgconn.PgError

	insertCredentials, err := r.db.PrepareContext(ctx, "INSERT INTO users (id, login, password) VALUES ($1, $2, $3);")
	if err != nil {
		return err
	}
	defer insertCredentials.Close()
	_, err = insertCredentials.ExecContext(ctx, id, login, hashedPassword)

	log.Println("1")
	if err != nil && errors.As(err, &pgErr) {
		log.Println("2")

		if pgErr.Code == uniqueViolationErr {
			log.Println("3")
			return entities.ErrLoginAlreadyInUse
		} else {
			log.Println("4")
			return err
		}
	}
	return nil
}

func (r *repository) GetCredentials(ctx context.Context, login string) (string, error) {
	var pass string

	selectPassForLogin, err := r.db.PrepareContext(
		ctx, "SELECT password FROM users WHERE login=$1;",
	)
	if err != nil {
		return "", err
	}
	defer selectPassForLogin.Close()

	err = selectPassForLogin.QueryRowContext(ctx, login).Scan(&pass)

	if err != nil {
		return "", err
	}
	return pass, nil
}

func (r *repository) Ping() error {
	ctx, cancel := context.WithTimeout(r.ctx, 1*time.Second)
	defer cancel()
	err := r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("PingContext failed: %w", err)
	}
	return nil
}

func (r *repository) Close() {
	r.db.Close()
}

func NewRepository(ctx context.Context, psqlConn string) *repository {
	db, err := sql.Open("pgx", psqlConn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	result, err := db.ExecContext(ctx, schema)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	return &repository{
		db:  db,
		ctx: ctx,
	}
}
