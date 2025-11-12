package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCategories(pool *pgxpool.Pool) ([]models.CategoryProduct, error) {
	var categories []models.CategoryProduct

	rows, err := pool.Query(context.Background(), "SELECT id, name FROM category_product")
	if err != nil {
		fmt.Println("Error: Failed to get category", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category models.CategoryProduct
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			fmt.Println("Error scanning category:", err)
			continue
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func CreateCategory(pool *pgxpool.Pool, input models.CategoryProduct) (models.CategoryProduct, error) {
	err := pool.QueryRow(context.Background(),
		"INSERT INTO category_product (name) VALUES ($1) RETURNING id", input.Name,
	).Scan(&input.Id)

	if err != nil {
		fmt.Println("Error: Failed type request")
	}

	return input, nil
}

func GetCategoryById(pool *pgxpool.Pool, id int) (models.CategoryProduct, error) {
	var category models.CategoryProduct
	err := pool.QueryRow(context.Background(),
		"SELECT id, name FROM category_product WHERE id=$1", id,
	).Scan(&category.Id, &category.Name)

	if err != nil {
		fmt.Println("Error: Category not found", err)
	}

	return category, nil
}

func EditCategory(pool *pgxpool.Pool, id int, input models.CategoryProduct) (models.CategoryProduct, error) {
	category, err := GetCategoryById(pool, id)
	if err != nil {
		fmt.Println("Error: Not Found Category", err)
	}

	if input.Name != "" {
		category.Name = input.Name
	}

	_, err = pool.Exec(context.Background(),
		"UPDATE category_product SET name=$1 WHERE id=$2",
		category.Name, id,
	)

	if err != nil {
		return category, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func DeleteCategory(pool *pgxpool.Pool, id int) error {
	res, err := pool.Exec(context.Background(),
		"DELETE FROM category_product WHERE id=$1", id,
	)

	if err != nil {
		fmt.Printf("failed to delete category: %w", err)
	}

	if res.RowsAffected() == 0 {
		fmt.Println("category not found")
	}

	return nil
}