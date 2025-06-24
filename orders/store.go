package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/logan2k02/ims/shared/utils"

	pb "github.com/logan2k02/ims/shared/protobuf"
)

type ordersStore struct {
	db *sql.DB
}

var (
	dbHost     = utils.GetEnv("DB_HOST", "localhost")
	dbPort     = utils.GetEnv("DB_PORT", "3306")
	dbUser     = utils.GetEnv("DB_USER", "admin")
	dbPassword = utils.GetEnv("DB_PASSWORD", "123456")
	dbName     = utils.GetEnv("DB_NAME", "ims_db")
)

func NewOrdersStore() (*ordersStore, error) {
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

	return &ordersStore{
		db: conn,
	}, nil
}

func (s *ordersStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *ordersStore) Init() error {
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
	CREATE TABLE IF NOT EXISTS orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		payment_reference VARCHAR(100),
		customer_name VARCHAR(255) NOT NULL,
		customer_contact VARCHAR(255) NOT NULL,
		status ENUM('pending', 'completed', 'cancelled') NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS order_items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		order_id INT NOT NULL,
		product_id INT NOT NULL,
		quantity INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE
	);
	`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

type jsonOrderItem struct {
	ProductId int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

func (s *ordersStore) rowToOrder(row *sql.Row) (*pb.Order, error) {
	var order pb.Order
	var itemsJSON []byte

	if err := row.Scan(&order.Id, &order.PaymentReference, &order.CustomerName, &order.CustomerContact, &order.Status, &order.CreatedAt, &itemsJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var items []jsonOrderItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
	}

	var orderItems []*pb.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, &pb.OrderItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	order.Items = orderItems
	return &order, nil
}

const SINGLE_ROW_ORDER_QUERY = `
SELECT 
	o.id,
	o.payment_reference,
	o.customer_name,
	o.customer_contact,
	o.status,
	o.created_at,
	COALESCE(
    JSON_ARRAYAGG(
			JSON_OBJECT(
				'product_id', oi.product_id,
				'quantity', oi.quantity
			)
		), JSON_ARRAY()
	) AS items
	FROM orders o
	LEFT JOIN order_items oi ON o.id = oi.order_id
	WHERE o.id = ?
	GROUP BY o.id;
	`

func (s *ordersStore) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("create order", "failed to rollback transaction: %v", err)
		}
	}()

	query := `
	INSERT INTO orders (payment_reference, customer_name, customer_contact, status)
	VALUES (?, ?, ?, 'pending')
	`

	result, err := tx.ExecContext(ctx, query, payload.PaymentReference, payload.CustomerName, payload.CustomerContact)
	if err != nil {
		return nil, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	query = `INSERT INTO order_items (order_id, product_id, quantity) VALUES `
	for i, item := range payload.Items {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("(%d, %d, %d)", orderID, item.ProductId, item.Quantity)
	}

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, SINGLE_ROW_ORDER_QUERY, orderID)

	order, err := s.rowToOrder(row)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *ordersStore) GetOrder(ctx context.Context, payload *pb.OrderIdRequest) (*pb.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("get order", "failed to rollback transaction: %v", err)
		}
	}()

	row := tx.QueryRowContext(ctx, SINGLE_ROW_ORDER_QUERY, payload.Id)

	order, err := s.rowToOrder(row)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *ordersStore) ListOrders(ctx context.Context, payload *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("list orders", "failed to rollback transaction: %v", err)
		}
	}()

	page := 1
	pageSize := 10
	if payload.Page > 0 {
		page = int(payload.Page)
	}

	if payload.PageSize > 0 {
		pageSize = int(payload.PageSize)
	}

	offset := (page - 1) * pageSize

	query := fmt.Sprintf(`
	SELECT 
	o.id,
	o.payment_reference,
	o.customer_name,
	o.customer_contact,
	o.status,
	o.created_at,
	COALESCE(
    JSON_ARRAYAGG(
			JSON_OBJECT(
				'product_id', oi.product_id,
				'quantity', oi.quantity
			)
		), JSON_ARRAY()
	) AS items
	FROM orders o
	LEFT JOIN order_items oi ON o.id = oi.order_id
	GROUP BY o.id
	ORDER BY o.created_at DESC
	LIMIT %d OFFSET %d;
	`, pageSize, offset)

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			Logger.LogError("list orders", "failed to close rows: %v", err)
		}
	}()

	var orders []*pb.Order
	for rows.Next() {
		var order pb.Order
		var itemsJson []byte

		if err := rows.Scan(&order.Id, &order.PaymentReference, &order.CustomerName, &order.CustomerContact, &order.Status, &order.CreatedAt, &itemsJson); err != nil {
			return nil, err
		}

		var items []jsonOrderItem
		if err := json.Unmarshal(itemsJson, &items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
		}

		var orderItems []*pb.OrderItem
		for _, item := range items {
			orderItems = append(orderItems, &pb.OrderItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			})
		}

		order.Items = orderItems
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	totalCountQuery := `SELECT COUNT(*) FROM orders;`
	var totalCount int64
	if err := tx.QueryRowContext(ctx, totalCountQuery).Scan(&totalCount); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pb.ListOrdersResponse{
		Orders:     orders,
		TotalCount: totalCount,
	}, nil
}

func (s *ordersStore) DeleteOrder(ctx context.Context, payload *pb.OrderIdRequest) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("delete order", "failed to rollback transaction: %v", err)
		}
	}()

	query := `DELETE FROM orders WHERE id = ?`
	if _, err := tx.ExecContext(ctx, query, payload.Id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *ordersStore) ChangeOrderStatus(ctx context.Context, payload *pb.ChangeOrderStatusRequest) (*pb.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			Logger.LogError("change order status", "failed to rollback transaction: %v", err)
		}
	}()

	query := `UPDATE orders SET status = ? WHERE id = ?`
	if _, err := tx.ExecContext(ctx, query, payload.Status, payload.Id); err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, SINGLE_ROW_ORDER_QUERY, payload.Id)
	order, err := s.rowToOrder(row)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}
