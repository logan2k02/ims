package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/logan2k02/ims/shared/utils"

	pb "github.com/logan2k02/ims/shared/protobuf"
)

type inventoryStore struct {
	db *sql.DB
}

var (
	dbHost     = utils.GetEnv("DB_HOST", "localhost")
	dbPort     = utils.GetEnv("DB_PORT", "3306")
	dbUser     = utils.GetEnv("DB_USER", "admin")
	dbPassword = utils.GetEnv("DB_PASSWORD", "123456")
	dbName     = utils.GetEnv("DB_NAME", "ims_db")
)

func NewInventoryStore() (*inventoryStore, error) {
	cfg := mysql.NewConfig()
	cfg.User = dbUser
	cfg.Passwd = dbPassword
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%s", dbHost, dbPort)
	cfg.DBName = dbName

	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	// ping
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &inventoryStore{
		db: conn,
	}, nil
}

func (s *inventoryStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *inventoryStore) Init() error {
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
	CREATE TABLE IF NOT EXISTS stock_movements (
		id INT AUTO_INCREMENT PRIMARY KEY,
		product_id INT NOT NULL,
		quantity_change INT NOT NULL,
    	type ENUM('purchase', 'supply', 'correction') NOT NULL,
		reference VARCHAR(100),
    	note TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE
	);
	`)
	if err != nil {
		return err
	}

	return tx.Commit()

}

type UpdateStockDto struct {
	ProductId int64
	Change    int64
	Reference string
	Note      string
	Type      string
}

func (s *inventoryStore) UpdateStockQuantity(ctx context.Context, payload *UpdateStockDto) (*pb.StockMovement, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("update stock quantity", "failed to rollback transaction: %v", err)
		}
	}()

	query := `
	UPDATE products
	SET stock_quantity = stock_quantity + ?
	WHERE id = ?
	`
	if payload.Type == "correction" {
		query = `
		UPDATE products
		SET stock_quantity = ?
		WHERE id = ?
		`
	}

	if _, err := tx.ExecContext(ctx, query, payload.Change, payload.ProductId); err != nil {
		return nil, err
	}

	quantityChange := payload.Change
	if quantityChange < 0 {
		quantityChange = -quantityChange
	}

	query = `
	INSERT INTO stock_movements (product_id, quantity_change, type, reference, note)
	VALUES (?,?,?,?,?)
	`

	result, err := tx.ExecContext(ctx, query, payload.ProductId, quantityChange, payload.Type, payload.Reference, payload.Note)
	if err != nil {
		return nil, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	var record pb.StockMovement

	row := tx.QueryRowContext(ctx, `SELECT id, product_id, quantity_change, type, reference, note, created_at FROM stock_movements WHERE id=?`, insertedId)
	if err := row.Scan(&record.Id, &record.ProductId, &record.Change, &record.Type, &record.Reference, &record.Note, &record.CreatedAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &record, nil
}

func (s *inventoryStore) ListStockMovements(ctx context.Context, productId int64) ([]*pb.StockMovement, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("get stock movements", "failed to rollback transaction: %v", err)
		}
	}()

	query := `
	SELECT id, product_id, quantity_change, type, reference, note, created_at FROM stock_movements
	`

	var args []any
	if productId > 0 {
		query += " WHERE id = ?"
		args = append(args, productId)
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*pb.StockMovement
	for rows.Next() {
		var record pb.StockMovement
		if err := rows.Scan(&record.Id, &record.ProductId, &record.Change, &record.Type, &record.Reference, &record.Note, &record.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
