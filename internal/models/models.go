package models

import "time"

type UnitOfMeasure string

const (
	UOMKg   UnitOfMeasure = "kg"
	UOMG    UnitOfMeasure = "g"
	UOML    UnitOfMeasure = "l"
	UOMMl   UnitOfMeasure = "ml"
	UOMUnit UnitOfMeasure = "unit"
	UOMBox  UnitOfMeasure = "box"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type RawMaterial struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Stock     float64       `json:"stock"`
	UOM       UnitOfMeasure `json:"uom"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Recipe struct {
	ProductID     int     `json:"product_id"`
	RawMaterialID int     `json:"raw_material_id"`
	Quantity      float64 `json:"quantity"`
}

type Order struct {
	ID          int         `json:"id"`
	TotalAmount float64     `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

type FinancialLedger struct {
	ID              int       `json:"id"`
	OrderID         *int      `json:"order_id,omitempty"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}
