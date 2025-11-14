package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCategories(pool *pgxpool.Pool) ([]models.CategoryProduct, error) {
	var categories []models.CategoryProduct

	rows, err := pool.Query(context.Background(), "SELECT id, name FROM category_products")
	if err != nil {
		return nil, fmt.Errorf("failed to get categories, %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category models.CategoryProduct
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func CreateCategory(pool *pgxpool.Pool, input models.CategoryProduct) (models.CategoryProduct, error) {
	err := pool.QueryRow(context.Background(),
		"INSERT INTO category_products (name) VALUES ($1) RETURNING id", input.Name,
	).Scan(&input.Id)

	if err != nil {
		return models.CategoryProduct{}, fmt.Errorf("failed to insert category, %w", err)
	}	

	return input, nil
}

func GetCategoryById(pool *pgxpool.Pool, id int) (models.CategoryProduct, error) {
	var category models.CategoryProduct
	err := pool.QueryRow(context.Background(),
		"SELECT id, name FROM category_products WHERE id=$1", id,
	).Scan(&category.Id, &category.Name)

	if err != nil {
		return  models.CategoryProduct{}, fmt.Errorf("error: Category not found, %w", err)
	}

	return category, nil
}

func EditCategory(pool *pgxpool.Pool, id int, input models.CategoryProduct) (models.CategoryProduct, error) {
	category, err := GetCategoryById(pool, id)
	if err != nil {
		return  models.CategoryProduct{}, fmt.Errorf("error: not found category, %w", err)
	}

	if input.Name != "" {
		category.Name = input.Name
	}

	_, err = pool.Exec(context.Background(),
		"UPDATE category_products SET name=$1 WHERE id=$2",
		category.Name, id,
	)

	if err != nil {
		return  models.CategoryProduct{}, fmt.Errorf("error: failed to update category, %w", err)
	}

	return category, nil
}

func DeleteCategory(pool *pgxpool.Pool, id int) error {
	_, err := pool.Exec(context.Background(),
		"DELETE FROM category_products WHERE id=$1", id,
	)

	if err != nil {
		return fmt.Errorf("error: not getting category, %w", err)
	}

	return nil
}