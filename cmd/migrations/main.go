package main

import (
	"back-end-coffeShop/controller"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

)

func main() {
	db := controller.ConnectDB()
	defer db.Close()

	migrationsDir := "../../lib/sql/migrations"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("gagal membaca folder migrations: %v", err)
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			filePath := filepath.Join(migrationsDir, f.Name())
			sqlBytes, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("gagal membaca file %s: %v", f.Name(), err)
			}

			query := string(sqlBytes)
			fmt.Printf("menjalankan migrasi: %s\n", f.Name())

			if _, err := db.Exec(context.Background(), query); err != nil {
				log.Fatalf("gagal %s: %v", f.Name(), err)
			}
		}
	}

	fmt.Println(" .up.sql berhasil dijalankan.")
}