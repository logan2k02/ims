package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/theshawa/imms/shared/protobuf"
	"github.com/theshawa/imms/shared/utils"
)

type productsStore struct {
	db *sql.DB
}

var (
	dbHost     = utils.GetEnv("DB_HOST", "localhost")
	dbPort     = utils.GetEnv("DB_PORT", "5432")
	dbUser     = utils.GetEnv("DB_USER", "admin")
	dbPassword = utils.GetEnv("DB_PASSWORD", "123456")
	dbName     = utils.GetEnv("DB_NAME", "imms_db")
)

func NewProductsStore() (*productsStore, error) {
	cfg := mysql.NewConfig()
	cfg.User = dbUser
	cfg.Passwd = dbPassword
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%s", dbHost, dbPort)
	cfg.DBName = dbName

	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db due to an error: %v", err)
	}

	// ping
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &productsStore{
		db: conn,
	}, nil
}

func (s *productsStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *productsStore) Init() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return fmt.Errorf("failed to initialize products table: %w", err)
	}

	return nil
}

func rowToProduct(row *sql.Row) (*protobuf.Product, error) {
	var product protobuf.Product
	if err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to scan product row: %w", err)
	}
	return &product, nil
}

func (s *productsStore) CreateProduct(ctx context.Context, payload *protobuf.CreateProductRequest) (*protobuf.Product, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	query := `
	INSERT INTO products (name, description, price)
	VALUES (?, ?, ?)
	`

	result, err := tx.ExecContext(ctx, query, payload.Name, payload.Description, payload.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to make query: %w", err)
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	row := tx.QueryRow(`SELECT id, name, description, price, created_at FROM products WHERE id = ?`, insertedId)

	insertedProduct, err := rowToProduct(row)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return insertedProduct, nil
}
