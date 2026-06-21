package pos

import (
	"context"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"integritypos/internal/models"
)

func ProcessOrder(ctx context.Context, db *pgxpool.Pool, items []models.OrderItem) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Subtotal
	}

	var orderID int
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (total_amount, status)
		VALUES ($1, 'completed')
		RETURNING id
	`, totalAmount).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	rawMaterialDeductions := make(map[int]float64)

	for _, item := range items {
		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
			VALUES ($1, $2, $3, $4, $5)
		`, orderID, item.ProductID, item.Quantity, item.UnitPrice, item.Subtotal)
		if err != nil {
			return fmt.Errorf("failed to insert order item for product_id %d: %w", item.ProductID, err)
		}

		rows, err := tx.Query(ctx, `
			SELECT raw_material_id, quantity
			FROM recipes
			WHERE product_id = $1
		`, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to query recipes for product_id %d: %w", item.ProductID, err)
		}

		for rows.Next() {
			var rmID int
			var qty float64
			if err := rows.Scan(&rmID, &qty); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan recipe for product_id %d: %w", item.ProductID, err)
			}
			rawMaterialDeductions[rmID] += qty * float64(item.Quantity)
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			return fmt.Errorf("rows iteration error for product_id %d recipes: %w", item.ProductID, err)
		}
	}

	var rmIDs []int
	for id := range rawMaterialDeductions {
		rmIDs = append(rmIDs, id)
	}
	sort.Ints(rmIDs)

	for _, rmID := range rmIDs {
		deduction := rawMaterialDeductions[rmID]

		var currentStock float64
		err := tx.QueryRow(ctx, `
			SELECT stock 
			FROM raw_materials 
			WHERE id = $1 
			FOR UPDATE
		`, rmID).Scan(&currentStock)
		if err != nil {
			if err == pgx.ErrNoRows {
				return fmt.Errorf("raw material %d not found for deduction", rmID)
			}
			return fmt.Errorf("failed to acquire lock for raw material %d: %w", rmID, err)
		}

		if currentStock < deduction {
			return fmt.Errorf("insufficient stock for raw material %d: current %.4f, required %.4f", rmID, currentStock, deduction)
		}

		_, err = tx.Exec(ctx, `
			UPDATE raw_materials 
			SET stock = stock - $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`, deduction, rmID)
		if err != nil {
			return fmt.Errorf("failed to deduct stock for raw material %d: %w", rmID, err)
		}
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO financial_ledger (order_id, transaction_type, amount, description)
		VALUES ($1, 'sale', $2, 'Order payment received')
	`, orderID, totalAmount)
	if err != nil {
		return fmt.Errorf("failed to record entry in financial ledger for order %d: %w", orderID, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
