package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Fprintln(os.Stderr, "warning: .env not found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is not set")
		os.Exit(1)
	}

	sql, err := os.ReadFile("seed.sql")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read seed.sql:", err)
		os.Exit(1)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect to database:", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	if _, err := conn.Exec(ctx, string(sql)); err != nil {
		fmt.Fprintln(os.Stderr, "failed to execute seed:", err)
		os.Exit(1)
	}

	fmt.Println("Seed executado com sucesso!")
	fmt.Println("  users:          8 registros (1 admin, 3 clients, 4 professionals)")
	fmt.Println("  professionals:  4 registros")
	fmt.Println("  availabilities: 6 registros")
	fmt.Println("  contracts:      5 registros (1 PENDING, 2 ACTIVE, 1 COMPLETED, 1 CANCELLED)")
	fmt.Println("  addresses:      6 registros (4 de usuários, 2 de contratos)")
	fmt.Println("  reviews:        1 registro (contrato COMPLETED)")
}
