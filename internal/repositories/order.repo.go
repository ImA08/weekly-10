package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"minitask1.go/internal/models"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) SetOrder(ctx context.Context, order models.OrderStruct, id int, invoice string) (pgconn.CommandTag, error) {
	tx, err := o.db.Begin(ctx)

	if err != nil {
		return pgconn.CommandTag{}, nil
	}
	query := `INSERT INTO transactions (user_id, schedule_id, fullname, email, phone_number, payment_method_id, transaction_date, invoice_number) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7) RETURNING id`
	var transactionsId int

	transactionValues := []any{id, order.Schedule_id, order.Fullname, order.Email, order.PhoneNumber, order.Payment_method_id, invoice}
	if err := tx.QueryRow(ctx, query, transactionValues...).Scan(&transactionsId); err != nil {
		return pgconn.CommandTag{}, err
	}

	detQuery := `INSERT INTO order_seats (transaction_id, seat_id) VALUES`

	orderValues := []any{transactionsId}
	for i, seat_id := range order.Seats {
		if i > 0 {
			detQuery += ", "
		}
		detQuery += fmt.Sprintf(` ($1, $%d)`, len(orderValues)+1)
		log.Println("[DEBUG] SEAT ID", seat_id)
		orderValues = append(orderValues, seat_id)
	}

	_, err = tx.Exec(ctx, detQuery, orderValues...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	queryTotal := `WITH sum_total AS (
		SELECT 
			s.price * COUNT(os.id) AS calculated_total
		FROM 
			schedules s 
		JOIN 
			transactions t ON s.id = t.schedule_id
		JOIN 
			order_seats os ON t.id = os.transaction_id
		WHERE 
			t.id = $1
		GROUP BY 
			s.price
		)
	UPDATE transactions
	SET total = (SELECT calculated_total FROM sum_total)
	WHERE id = $1`

	cmd, err := tx.Exec(ctx, queryTotal, transactionsId)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil

}

func (o OrderRepository) GetOrder(ctx context.Context, transactionId int) (*models.OrderStruct, error) {
	query := `SELECT
t.user_id,
t.id,
t.fullname,
t.email,
t.schedule_id,
array_agg(s.name),
t.transaction_date,
t.status_paid,
t.total,
t.updated_at,
t.invoice_number,
t.payment_method_id,
t.phone_number
FROM transactions t
LEFT JOIN order_seats os ON os.transaction_id = t.id
JOIN seats s ON s.id = os.seat_id
WHERE t.id = $1
GROUP BY t.id`

	var order models.OrderStruct

	err := o.db.QueryRow(ctx, query, transactionId).Scan(
		&order.User_id,
		&order.Id,
		&order.Fullname,
		&order.Email,
		&order.Schedule_id,
		&order.Seats,
		&order.Transaction_date,
		&order.Status_paid,
		&order.Total,
		&order.Updated_at,
		&order.Invoice_number,
		&order.Payment_method_id,
		&order.PhoneNumber,
	)

	if err != nil {
		return nil, err
	}

	return &order, err
}

func (o *OrderRepository) UpdateOrder(ctx context.Context, transactionId int) (models.OrderStruct, error) {
	query := `UPDATE transactions SET status_paid = true WHERE id = $1`

	values := []any{transactionId}

	var result models.OrderStruct

	if err := o.db.QueryRow(ctx, query, values...).Scan(&result.Status_paid); err != nil {
		return models.OrderStruct{}, err
	}

	return result, nil
}
