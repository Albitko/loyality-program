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
	  	"order_number" text primary key unique,
	  	user_id text not null references users(id),
	    status text not null,
	    accrual float not null default 0,
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
	updateOrder, err := r.db.PrepareContext(
		ctx, "UPDATE orders SET status=$1, accrual=$2, current=$2 WHERE order_number=$3;",
	)
	log.Println(order)
	defer updateOrder.Close()
	if err != nil {
		return err
	}
	_, err = updateOrder.ExecContext(ctx, order.Status, order.Accrual, order.OrderID)
	if err != nil {
		return err
	}
	return nil
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
	var userID string

	selectUserIDForOrder, err := r.db.PrepareContext(
		ctx, "SELECT user_id FROM orders WHERE order_number=$1;",
	)
	defer selectUserIDForOrder.Close()
	if err != nil {
		return "", err
	}
	err = selectUserIDForOrder.QueryRowContext(ctx, order).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", entities.ErrNoOrderForUser
		} else {
			return "", err
		}
	}
	return userID, nil
}

func (r *repository) CreateOrder(ctx context.Context, order entities.Order, userID string) error {
	now := time.Now()
	uploadedAt := now.Format(time.RFC3339)

	createOrder, err := r.db.PrepareContext(
		ctx, "INSERT INTO orders (order_number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4);",
	)
	defer createOrder.Close()

	if err != nil {
		return err
	}
	_, err = createOrder.ExecContext(ctx, order.OrderID, userID, order.Status, uploadedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetOrdersForUser(ctx context.Context, userID string) ([]entities.OrderWithTime, error) {
	var orders []entities.OrderWithTime
	var order entities.OrderWithTime

	selectOrdersForUser, err := r.db.PrepareContext(
		ctx,
		"SELECT order_number, status, accrual, uploaded_at FROM orders WHERE user_id=$1 ORDER BY uploaded_at;",
	)
	defer selectOrdersForUser.Close()
	if err != nil {
		return orders, err
	}
	row, err := selectOrdersForUser.QueryContext(ctx, userID)
	defer row.Close()
	if err != nil {
		return orders, err
	}

	if err = row.Err(); err != nil {
		return orders, err
	}
	for row.Next() {
		err := row.Scan(&order.OrderID, &order.Status, &order.Accrual, &order.UpdatedAt)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return orders, entities.ErrNoOrderForUser
	}
	return orders, nil
}

func (r *repository) Register(ctx context.Context, id, login, hashedPassword string) error {
	var pgErr *pgconn.PgError

	insertCredentials, err := r.db.PrepareContext(
		ctx, "INSERT INTO users (id, login, password) VALUES ($1, $2, $3);",
	)
	if err != nil {
		return err
	}
	defer insertCredentials.Close()
	_, err = insertCredentials.ExecContext(ctx, id, login, hashedPassword)

	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == uniqueViolationErr {
			return entities.ErrLoginAlreadyInUse
		} else {
			return err
		}
	}
	return nil
}

func (r *repository) GetCredentials(ctx context.Context, login string) (entities.User, error) {
	var user entities.User
	var id string
	var hashedPassword string

	selectPassForLogin, err := r.db.PrepareContext(
		ctx, "SELECT id, password FROM users WHERE login=$1;",
	)
	if err != nil {
		return user, err
	}
	defer selectPassForLogin.Close()

	err = selectPassForLogin.QueryRowContext(ctx, login).Scan(&id, &hashedPassword)
	if err != nil {
		return user, err
	}
	user.ID = id
	user.Login = login
	user.Password = hashedPassword
	return user, nil
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
