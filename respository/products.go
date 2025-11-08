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

func GetById(ctx *gin.Context, conn *pgx.Conn) (models.Product, error) {
	id := ctx.Param("id")

	var product models.Product

	err := conn.QueryRow(context.Background(), `
		SELECT id, name, price, description, productsize, stock, isflashsale, tempelatur, category_productid, created_at, updated_at
		FROM product WHERE id = $1
	`, id).Scan(
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

	return product, err
}

func Edit(conn *pgx.Conn, ctx *gin.Context) (models.Product, error) {
	id := ctx.Param("id")

	oldProduct, err := GetById(ctx, conn)

	if err != nil {
		return models.Product{}, fmt.Errorf("product not found")
	}

	var newProduct models.Product

	err = ctx.BindJSON(&newProduct)

	if err != nil {
		fmt.Println("Error : Failed type request much json type")
	}

	if newProduct.Name == "" {
		newProduct.Name = oldProduct.Name
	}
	if newProduct.Price == 0 {
		newProduct.Price = oldProduct.Price
	}
	if newProduct.Description == "" {
		newProduct.Description = oldProduct.Description
	}
	if newProduct.Productsize == "" {
		newProduct.Productsize = oldProduct.Productsize
	}
	if newProduct.Stock == 0 {
		newProduct.Stock = oldProduct.Stock
	}
	if newProduct.Isflashsale == nil {
		newProduct.Isflashsale = oldProduct.Isflashsale
	}
	if newProduct.Tempelatur == "" {
		newProduct.Tempelatur = oldProduct.Tempelatur
	}
	if newProduct.Category_productid == 0 {
		newProduct.Category_productid = oldProduct.Category_productid
	}

	_, err = conn.Exec(context.Background(), "UPDATE product SET name=$1, price=$2, description=$3, productsize=$4, stock=$5, isflashsale=$6, tempelatur=$7, category_productid=$8, updated_at=NOW() WHERE id = $9", newProduct.Name, newProduct.Price, newProduct.Description, newProduct.Productsize,
		newProduct.Stock, *newProduct.Isflashsale, newProduct.Tempelatur, newProduct.Category_productid, id)

	return newProduct, err
}
