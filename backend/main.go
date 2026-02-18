package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/api"
	"github.com/durianpay/fullstack-boilerplate/internal/config"
	ah "github.com/durianpay/fullstack-boilerplate/internal/module/auth/handler"
	ar "github.com/durianpay/fullstack-boilerplate/internal/module/auth/repository"
	au "github.com/durianpay/fullstack-boilerplate/internal/module/auth/usecase"
	ph "github.com/durianpay/fullstack-boilerplate/internal/module/payment/handler"
	pr "github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
	pu "github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	srv "github.com/durianpay/fullstack-boilerplate/internal/service/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	db, err := sql.Open("sqlite3", "dashboard.db?_foreign_keys=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	JwtExpiredDuration, err := time.ParseDuration(config.JwtExpired)
	if err != nil {
		panic(err)
	}

	userRepo := ar.NewUserRepo(db)
	authUC := au.NewAuthUsecase(userRepo, config.JwtSecret, JwtExpiredDuration)
	authH := ah.NewAuthHandler(authUC)

	paymentRepo := pr.NewPaymentRepository(db)
	paymentUC := pu.NewPaymentUsecase(paymentRepo)
	paymentH := ph.NewPaymentHandler(paymentUC)

	apiHandler := &api.APIHandler{
		Auth:    authH,
		Payment: paymentH,
	}

	server := srv.NewServer(apiHandler, config.OpenapiYamlLocation, config.JwtSecret)

	addr := config.HttpAddress
	log.Printf("starting server on %s", addr)
	server.Start(addr)
}

func initDB(db *sql.DB) error {
	// create tables if not exists
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  email TEXT NOT NULL UNIQUE,
		  password_hash TEXT NOT NULL,
		  role TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS merchants (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  name TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS orders (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  merchant_id INTEGER NOT NULL,
		  total INTEGER NOT NULL,
		  status TEXT NOT NULL,
		  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS payments (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  order_id INTEGER NOT NULL,
		  merchant_id INTEGER NOT NULL,
		  merchant_name TEXT NOT NULL,
		  amount INTEGER NOT NULL,
		  status TEXT NOT NULL,
		  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_merchant_id ON payments(merchant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_merchant_id ON orders(merchant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	// seed admin user if not exists
	var cnt int
	row := db.QueryRow("SELECT COUNT(1) FROM users")
	if err := row.Scan(&cnt); err != nil {
		return err
	}
	if cnt == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?)", "cs@test.com", string(hash), "cs"); err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?)", "operation@test.com", string(hash), "operation"); err != nil {
			return err
		}
	}

	// seed merchant if not exists
	var cntMerchant int
	rowMerchant := db.QueryRow("SELECT COUNT(1) FROM merchants")
	if err := rowMerchant.Scan(&cntMerchant); err != nil {
		return err
	}
	if cntMerchant == 0 {
		if err := seedMerchants(db); err != nil {
			return err
		}
	}

	var cntOrder int
	rowOrder := db.QueryRow("SELECT COUNT(1) FROM orders")
	if err := rowOrder.Scan(&cntOrder); err != nil {
		return err
	}
	if cntOrder == 0 {
		if err := seedOrders(db); err != nil {
			return err
		}
	}

	// seed payment if not exists
	var cntPayment int
	rowPayment := db.QueryRow("SELECT COUNT(1) FROM payments")
	if err := rowPayment.Scan(&cntPayment); err != nil {
		return err
	}
	if cntPayment == 0 {
		if err := seedPayments(db); err != nil {
			return err
		}
	}

	const dbLifetime = time.Minute * 5
	db.SetConnMaxLifetime(dbLifetime)
	return nil
}

func seedMerchants(db *sql.DB) error {
	var cntMerchant int
	err := db.QueryRow("SELECT COUNT(1) FROM merchants").Scan(&cntMerchant)
	if err != nil {
		return err
	}
	if cntMerchant > 0 {
		return nil
	}

	merchants := []string{
		"Tech Solutions Inc",
		"Global Traders Ltd",
		"Prime Services Co",
		"Digital Ventures",
		"Elite Products Corp",
	}

	for _, name := range merchants {
		if _, err := db.Exec("INSERT INTO merchants(name) VALUES (?)", name); err != nil {
			return err
		}
	}
	return nil
}

func seedOrders(db *sql.DB) error {
	var cntOrder int
	err := db.QueryRow("SELECT COUNT(1) FROM orders").Scan(&cntOrder)
	if err != nil {
		return err
	}
	if cntOrder > 0 {
		return nil
	}

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()
	statuses := []string{"draft", "pending", "paid", "processing", "shipped", "completed", "cancelled"}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 100; i++ {
		merchantID := rng.Intn(5) + 1
		totalAmount := rng.Intn(900000) + 10000
		status := statuses[rng.Intn(len(statuses))]
		randomTime := startDate.Add(
			time.Duration(rng.Int63n(endDate.Unix()-startDate.Unix())) * time.Second,
		)
		if _, err := db.Exec("INSERT INTO orders(merchant_id, total, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", merchantID, totalAmount, status, randomTime, randomTime); err != nil {
			return err
		}
	}
	return nil
}

func seedPayments(db *sql.DB) error {
	var cntPayment int
	err := db.QueryRow("SELECT COUNT(1) FROM payments").Scan(&cntPayment)
	if err != nil {
		return err
	}
	if cntPayment > 0 {
		return nil
	}

	// Get all orders with merchant info
	orders, err := fetchOrdersForPayments(db)
	if err != nil {
		return err
	}

	if len(orders) == 0 {
		return fmt.Errorf("no orders found, seed orders first")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create payments for each order
	for _, order := range orders {
		if err := createPaymentsForOrder(db, order, rng); err != nil {
			return err
		}
	}

	return nil
}

// OrderInfo holds order data needed for creating payments
type OrderInfo struct {
	ID           int
	MerchantID   int
	MerchantName string
	Total        int
	Status       string
	CreatedAt    time.Time
}

// fetchOrdersForPayments retrieves all orders with merchant info
func fetchOrdersForPayments(db *sql.DB) ([]OrderInfo, error) {
	rows, err := db.Query(`
		SELECT 
			o.id, 
			o.merchant_id, 
			m.name, 
			o.total, 
			o.status,
			o.created_at
		FROM orders o 
		JOIN merchants m ON o.merchant_id = m.id
		ORDER BY o.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []OrderInfo
	for rows.Next() {
		var order OrderInfo
		if err := rows.Scan(
			&order.ID,
			&order.MerchantID,
			&order.MerchantName,
			&order.Total,
			&order.Status,
			&order.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// createPaymentsForOrder creates payment(s) for a single order
func createPaymentsForOrder(db *sql.DB, order OrderInfo, rng *rand.Rand) error {
	if shouldSkipPayment(order.Status) {
		return nil
	}

	numPayments := determinePaymentAttempts(rng)

	for attempt := 1; attempt <= numPayments; attempt++ {
		payment := generatePayment(order, attempt, numPayments, rng)

		if err := insertPayment(db, payment); err != nil {
			return err
		}
	}

	return nil
}

// shouldSkipPayment checks if order should have payments
func shouldSkipPayment(orderStatus string) bool {
	return orderStatus == "draft" || orderStatus == "cancelled"
}

// determinePaymentAttempts decides how many payment attempts to create
func determinePaymentAttempts(rng *rand.Rand) int {
	roll := rng.Intn(100)

	if roll < 5 {
		return 3 // 5% - multiple retries
	} else if roll < 15 {
		return 2 // 10% - one retry
	}
	return 1 // 85% - single payment
}

// PaymentData holds data for a single payment
type PaymentData struct {
	OrderID      int
	MerchantID   int
	MerchantName string
	Amount       int
	Status       string
	CreatedAt    time.Time
}

// generatePayment creates payment data based on order and attempt number
func generatePayment(order OrderInfo, attempt, totalAttempts int, rng *rand.Rand) PaymentData {
	// Calculate payment time (0-60 mins after order, plus retry delays)
	paymentDelay := time.Duration(rng.Intn(3600)) * time.Second
	if attempt > 1 {
		paymentDelay += time.Duration(attempt*30) * time.Minute
	}
	paymentTime := order.CreatedAt.Add(paymentDelay)

	status := determinePaymentStatus(order.Status, attempt, totalAttempts, rng)

	amount := order.Total
	if attempt == 1 && totalAttempts == 1 && rng.Intn(10) < 1 {
		amount = order.Total / 2 // 10% chance of partial payment
	}

	return PaymentData{
		OrderID:      order.ID,
		MerchantID:   order.MerchantID,
		MerchantName: order.MerchantName,
		Amount:       amount,
		Status:       status,
		CreatedAt:    paymentTime,
	}
}

// determinePaymentStatus decides payment status based on order status and attempt
func determinePaymentStatus(orderStatus string, attempt, totalAttempts int, rng *rand.Rand) string {
	// Earlier attempts always fail
	if attempt < totalAttempts {
		return "failed"
	}

	// Last attempt - depends on order status
	completedOrderStatuses := []string{"paid", "processing", "shipped", "completed"}
	for _, status := range completedOrderStatuses {
		if orderStatus == status {
			return "completed"
		}
	}

	// Order is pending - randomize final payment status
	roll := rng.Intn(10)
	if roll < 7 {
		return "completed"
	} else if roll < 9 {
		return "processing"
	}
	return "failed"
}

// insertPayment inserts a payment record into the database
func insertPayment(db *sql.DB, payment PaymentData) error {
	_, err := db.Exec(`
		INSERT INTO payments 
		(order_id, merchant_id, merchant_name, amount, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, payment.OrderID, payment.MerchantID, payment.MerchantName,
		payment.Amount, payment.Status, payment.CreatedAt, payment.CreatedAt)

	return err
}
