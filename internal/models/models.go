package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	CreateAt       time.Time `json:"-"`
	UpdateAt       time.Time `json:"-"`
}

type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customer_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreateAt      time.Time `json:"-"`
	UpdateAt      time.Time `json:"-"`
}

type Status struct {
	ID       int       `json:"id"`
	Name     int       `json:"name"`
	CreateAt time.Time `json:"-"`
	UpdateAt time.Time `json:"-"`
}

type TransactionStatus struct {
	ID       int       `json:"id"`
	Name     int       `json:"name"`
	CreateAt time.Time `json:"-"`
	UpdateAt time.Time `json:"-"`
}

type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreateAt            time.Time `json:"-"`
	UpdateAt            time.Time `json:"-"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreateAt  time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreateAt  time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		select 
			id, 
			name, 
			description, 
			inventory_level, 
			price, 
			image,
			created_at,
			updated_at
		from 
			widgets 
		where 
			id = ?`

	var widget Widget
	row := m.DB.QueryRowContext(ctx, stmt, id)

	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.CreateAt,
		&widget.UpdateAt,
	)
	if err != nil {
		return widget, err
	}

	return widget, nil
}

func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into transactions 
		(amount, currency, last_four, expiry_month, expiry_year, payment_intent, payment_method,
			bank_return_code, transaction_status_id, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(
		ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into orders
		(widget_id, transaction_id, customer_id, status_id, quantity, amount, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(
		ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.CustomerID,
		order.StatusID,
		order.Quantity,
		order.Amount,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *DBModel) InsertCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into customers (first_name, last_name, email, created_at, updated_at)
		values (?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(
		ctx, stmt,
		c.FirstName,
		c.LastName,
		c.Email,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}
