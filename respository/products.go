package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetProducts(conn *pgx.Conn) ([]models.Product, error) {
	var dataProduct []models.Product

	rows, err := conn.Query(context.Background(), "SELECT id, name, price, description, productsize, stock, isflashsale, tempelatur, category_productid, created_at, updated_at FROM product")

	if err != nil {
		fmt.Println("Error: Failed get data product")
	}

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Productsize,
			&product.Stock,
			&product.Isflashsale,
			&product.Tempelatur,
			&product.Category_productid,
			&product.Created_at,
			&product.Updated_at,
		)
		if err != nil {
			fmt.Println("Error scanning product:", err)
		}
		dataProduct = append(dataProduct, product)
	}

	return dataProduct, nil
}

func Create(ctx *gin.Context, conn *pgx.Conn) models.Product {
	var input models.Product

	err := ctx.BindJSON(&input)

	if err != nil {
		fmt.Println("Error: Invalid type much json")
	}

	now := time.Now()

	_, err = conn.Exec(context.Background(), "INSERT INTO product (name, price, description, productsize, stock, isflashsale, tempelatur, category_productid, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", input.Name, input.Price, input.Description, input.Productsize, input.Stock, input.Isflashsale, input.Tempelatur, input.Category_productid, now, now)

	if err != nil {
		fmt.Println("Error insert product:", err)
	}

	input.Created_at = now
	input.Updated_at = now

	return input
}
