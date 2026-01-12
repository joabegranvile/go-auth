// @title Auth-Go API
// @version 1.0
// @description API de autenticação
// @host auth.127.0.0.1.lvh.me
// @BasePath /
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joabegranvile/auth-go/docs"
	"github.com/joabegranvile/auth-go/internal/auth"
	"github.com/joabegranvile/auth-go/internal/db"
	"github.com/joabegranvile/auth-go/internal/rbac"
	"github.com/joabegranvile/auth-go/internal/server"
	_ "github.com/lib/pq" // <--- IMPORTANTE: Driver do Postgres adicionado
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPass := getSecret("pg_password")

	// CORREÇÃO AQUI: user=%s (faltava o igual)
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName)

	pg, err := db.NewPostgres(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	if err := pg.Migrate(); err != nil {
		log.Fatal(err)
	}

	authSvc := auth.New("super-secret")
	rbacSvc := rbac.New()
	srv := server.New(authSvc, rbacSvc)
	srv.Routes()

	log.Println("HTTP :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
