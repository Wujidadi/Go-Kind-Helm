package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "✅ Go server is running.")
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		host := "postgres-go-helm-postgresql.postgres-go-helm.svc.cluster.local"
		port := "5432"
		dbname := "goappdb"
		user := "devuser"
		password := "devpassword"

		query := r.URL.Query()
		if _, ok := query["prod"]; ok {
			user = "produser"
			password = "strongprodpassword"
		}

		dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			http.Error(w, fmt.Sprintf("❌ 連線失敗（開啟時）: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			http.Error(w, fmt.Sprintf("❌ 連線失敗（Ping 時）: %v", err), http.StatusInternalServerError)
			return
		}

		// 試著查詢 PostgreSQL 版本
		var version string
		err = db.QueryRow("SELECT version()").Scan(&version)
		if err != nil {
			http.Error(w, fmt.Sprintf("✅ 成功連上 DB，但查詢版本失敗: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "✅ 成功連線到 PostgreSQL！\nPostgreSQL 版本：%s\n", version)
	})

	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		password := "dev-redis-password"

		query := r.URL.Query()
		if _, ok := query["prod"]; ok {
			password = "prod-redis-password"
		}

		rdb := redis.NewClient(&redis.Options{
			Addr:     "redis-go-helm-master.redis-go-helm.svc.cluster.local:6379",
			Password: password,
		})

		err := rdb.Set(ctx, "test-key", "HaHaHa Redis! I am Taras!", 0).Err()
		if err != nil {
			http.Error(w, fmt.Sprintf("❌ Redis SET error: %v", err), http.StatusInternalServerError)
			return
		}

		val, err := rdb.Get(ctx, "test-key").Result()
		if err != nil {
			http.Error(w, fmt.Sprintf("❌ Redis GET error: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "✅ Redis 測試成功！取得值: %s\n", val)
	})

	port := "2222"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Starting HTTP on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
