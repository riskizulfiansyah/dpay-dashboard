package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/api"
	"github.com/durianpay/fullstack-boilerplate/internal/config"
	ah "github.com/durianpay/fullstack-boilerplate/internal/module/auth/handler"
	ar "github.com/durianpay/fullstack-boilerplate/internal/module/auth/repository"
	au "github.com/durianpay/fullstack-boilerplate/internal/module/auth/usecase"
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

	apiHandler := &api.APIHandler{
		Auth: authH,
	}

	server := srv.NewServer(apiHandler, config.OpenapiYamlLocation)

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
		  name TEXT NOT NULL,
		  email TEXT NOT NULL UNIQUE,
		  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS orders (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  merchant_id INTEGER NOT NULL,
		  subtotal INTEGER NOT NULL,
		  status TEXT NOT NULL,
		  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS payments (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  order_id INTEGER NOT NULL,
		  amount INTEGER NOT NULL,
		  status TEXT NOT NULL,
		  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_merchant_id ON orders(merchant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);`,
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
		if _, err := db.Exec("INSERT INTO merchants(name, email) VALUES (?, ?)", "merchant1", "merchant1@merchant.com"); err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO merchants(name, email) VALUES (?, ?)", "merchant2", "merchant2@merchant.com"); err != nil {
			return err
		}
	}

	// seed order if not exists
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

	statuses := []string{"pending", "paid", "shipped", "completed", "cancelled"}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		// random timestamp between startDate and endDate
		randomTime := startDate.Add(
			time.Duration(rng.Int63n(endDate.Unix()-startDate.Unix())) * time.Second,
		)

		subtotal := rng.Intn(900000) + 10000 // 10,000 - 910,000
		status := statuses[rng.Intn(len(statuses))]
		merchantID := rng.Intn(2) + 1 // 1–2

		_, err := db.Exec(`
			INSERT INTO orders (subtotal, merchant_id, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
		`, subtotal, merchantID, status, randomTime, randomTime)
		if err != nil {
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

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()

	statuses := []string{"completed", "processing", "failed"}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 120; i++ {
		randomTime := startDate.Add(
			time.Duration(rng.Int63n(endDate.Unix()-startDate.Unix())) * time.Second,
		)

		orderID := rng.Intn(100) + 1 // 1–100
		amount := rng.Intn(900000) + 10000
		status := statuses[rng.Intn(len(statuses))]

		_, err := db.Exec(`
			INSERT INTO payments 
			(order_id, amount, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
		`, orderID, amount, status, randomTime, randomTime)
		if err != nil {
			return err
		}
	}

	return nil
}
