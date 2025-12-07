package main

import (
	"back-end-coffeShop/internal/controller"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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

	var downFiles []string

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".down.sql") {
			downFiles = append(downFiles, f.Name())
		}
	}

	if len(downFiles) == 0 {
		log.Println("Tidak ada file .down.sql ditemukan.")
		return
	}

	sort.Sort(sort.Reverse(sort.StringSlice(downFiles)))

	for _, fileName := range downFiles {
		filePath := filepath.Join(migrationsDir, fileName)
		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("gagal membaca file %s: %v", fileName, err)
		}

		query := string(sqlBytes)
		fmt.Printf("menjalankan rollback: %s\n", fileName)

		if _, err := db.Exec(context.Background(), query); err != nil {
			log.Fatalf("gagal rollback %s: %v", fileName, err)
		}
	}

	fmt.Println("Semua migration .down.sql berhasil dijalankan.")
}
