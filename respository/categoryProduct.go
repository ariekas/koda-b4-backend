package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCategorys(pool *pgxpool.Pool) ([]models.CategoryProduct, error){
	var dataCategory []models.CategoryProduct

	rows, err := pool.Query(context.Background(), "SELECT id, name FROM category_product")
	if err != nil {
		fmt.Println("Error: Failed to get category", err)
	}

	for rows.Next(){
		var category models.CategoryProduct
		err := rows.Scan(
			&category.Id,
			&category.Name,
		)
		if err != nil {
			fmt.Println("Error: Failed to scanning category", err)
		}

		dataCategory = append(dataCategory, category)
	}
	return dataCategory, nil
}

func CreateCategory(pool *pgxpool.Pool, ctx *gin.Context) (models.CategoryProduct, error) {
	var input models.CategoryProduct

	err := ctx.BindJSON(&input)
	if err != nil {
		fmt.Println("Error: Failed type request")
	}

	_, err = pool.Exec(context.Background(), "INSERT INTO category_product (name) VALUES ($1)", input.Name)

	return input, err
}

func GetCategoryById(pool *pgxpool.Pool, ctx *gin.Context) (models.CategoryProduct, error) {
	id := ctx.Param("id")

	var category models.CategoryProduct

	err := pool.QueryRow(context.Background(), "SELECT id, name FROM category_product WHERE id = $1", id).Scan(&category.Id, &category.Name)

	return category, err

}

func EditCategory(pool *pgxpool.Pool, ctx *gin.Context) (models.CategoryProduct, error) {
	id := ctx.Param("id")

	oldCategory, err := GetCategoryById(pool, ctx)
	if err != nil {
		fmt.Println("Error: Not Found Category", err)
	}

	var newCategory models.CategoryProduct

	err = ctx.BindJSON(&newCategory)

	if err != nil {
		fmt.Println("Error : Failed type request much json type")
	}

	if newCategory.Name == "" {
		newCategory.Name = oldCategory.Name
	}

	_, err = pool.Exec(context.Background(), "UPDATE category_product SET name=$1 WHERE id = $2", newCategory.Name, id)

	return newCategory, err
}

func DeleteCategory(pool *pgxpool.Pool, ctx *gin.Context) error {
	id := ctx.Param("id")
	_, err := pool.Exec(context.Background(), "DELETE FROM category_product WHERE id =$1", id)

	return err
}