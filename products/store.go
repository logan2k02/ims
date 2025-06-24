package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	pb "github.com/logan2k02/ims/shared/protobuf"
	"github.com/logan2k02/ims/shared/utils"
)

type productsStore struct {
	db *sql.DB
}

var (
	dbHost     = utils.GetEnv("DB_HOST", "localhost")
	dbPort     = utils.GetEnv("DB_PORT", "3306")
	dbUser     = utils.GetEnv("DB_USER", "admin")
	dbPassword = utils.GetEnv("DB_PASSWORD", "123456")
	dbName     = utils.GetEnv("DB_NAME", "ims_db")
)

func NewProductsStore() (*productsStore, error) {
	cfg := mysql.NewConfig()
	cfg.User = dbUser
	cfg.Passwd = dbPassword
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%s", dbHost, dbPort)
	cfg.DBName = dbName
	cfg.Timeout = 50 * time.Second

	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	// ping
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(0)
	conn.SetMaxOpenConns(500)

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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("init products store", "failed to rollback transaction: %v", err)
		}
	}()

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS products (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		sku VARCHAR(255) UNIQUE, 
		description TEXT, 
		price DECIMAL(10, 2) NOT NULL,
		reorder_level INT DEFAULT 0,
    	reorder_quantity INT DEFAULT 0,
		stock_quantity INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func rowToProduct(row *sql.Row) (*pb.Product, error) {
	var product pb.Product
	if err := row.Scan(&product.Id, &product.Name, &product.Sku, &product.Description, &product.Price, &product.ReorderLevel, &product.ReorderQuantity, &product.StockQuantity, &product.CreatedAt); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *productsStore) CreateProduct(ctx context.Context, payload *pb.CreateProductRequest) (*pb.Product, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("create product", "failed to rollback transaction: %v", err)
		}
	}()

	query := `
	INSERT INTO products (name, sku, description, price, reorder_level, reorder_quantity, stock_quantity)
	VALUES (?,?,?,?,?,?,?)
	`

	result, err := tx.ExecContext(ctx, query, payload.Name, payload.Sku, payload.Description, payload.Price, payload.ReorderLevel, payload.ReorderQuantity, payload.InitialQuantity)
	if err != nil {
		return nil, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, `SELECT id, name, sku, description, price, reorder_level, reorder_quantity, stock_quantity, created_at FROM products WHERE id = ?`, insertedId)

	insertedProduct, err := rowToProduct(row)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return insertedProduct, nil
}

func (s *productsStore) GetProducts(ctx context.Context, ids []int64) ([]*pb.Product, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("get products", "failed to rollback transaction: %v", err)
		}
	}()

	query := "SELECT id,name,sku,description,price, reorder_level, reorder_quantity, stock_quantity,created_at FROM products"
	args := make([]any, len(ids))
	if len(ids) > 0 {
		placeholders := make([]string, len(ids))
		for i, id := range ids {
			placeholders[i] = "?"
			args[i] = id
		}
		query += fmt.Sprintf(" WHERE id IN (%s)", strings.Join(placeholders, ", "))
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var product pb.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Sku, &product.Description, &product.Price, &product.ReorderLevel, &product.ReorderQuantity, &product.StockQuantity, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productsStore) DeleteProduct(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("delete product", "failed to rollback transaction: %v", err)
		}
	}()

	query := "DELETE FROM products WHERE id = ?"

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *productsStore) UpdateProduct(ctx context.Context, payload *pb.UpdateProductRequest) (*pb.Product, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("update product", "failed to rollback transaction: %v", err)
		}
	}()

	query := `
	UPDATE products
	SET name = ?,sku = ?, description = ?, price = ?, reorder_level=?, reorder_quantity=?
	WHERE id = ?
	`

	_, err = tx.ExecContext(ctx, query, payload.Name, payload.Sku, payload.Description, payload.Price, payload.ReorderLevel, payload.ReorderQuantity, payload.Id)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, `SELECT id, name, sku, description, price,reorder_level,reorder_quantity,stock_quantity, created_at FROM products WHERE id = ?`, payload.Id)

	updatedProduct, err := rowToProduct(row)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updatedProduct, nil
}
