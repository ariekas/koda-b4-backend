package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaginationResponse struct {
	Data       []models.Product  `json:"data"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int               `json:"total"`
	TotalPages int               `json:"total_pages"`
	Links      map[string]string `json:"links"`
}

func GetProducts(pool *pgxpool.Pool, page int) (PaginationResponse, error) {
	var dataProduct []models.Product
	limit := 50
	offset := (page - 1) * limit

	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		fmt.Println("Error counting products:", err)
	}

	rows, err := pool.Query(context.Background(), `
		SELECT p.id, p.name, p.price, p.description, p.product_size, p.stock, 
		       p.isFlashSale, p.isFavorite_product, p.temperature, p.category_product_id, 
		       COALESCE(ip.image, '') AS image,
		       p.created_at, p.updated_at
		FROM products p
		LEFT JOIN image_products ip ON ip.product_id = p.id
		ORDER BY p.id
		OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		fmt.Println("Error: Failed get data product", err)
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
			&product.IsFavorite_product,
			&product.Tempelatur,
			&product.Category_productid,
			&product.Image,
			&product.Created_at,
			&product.Updated_at,
		)
		if err != nil {
			fmt.Println("Error scanning product:", err)
		}
		dataProduct = append(dataProduct, product)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	links := make(map[string]string)

	if page > 1 {
		links["prev"] = fmt.Sprintf("/products?page=%d", page-1)
	} else {
		links["prev"] = "null"
	}

	if page < totalPages {
		links["next"] = fmt.Sprintf("/products?page=%d", page+1)
	}

	return PaginationResponse{
		Data:       dataProduct,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}, nil
}

func Create(ctx *gin.Context, pool *pgxpool.Pool) models.Product {
	var input models.Product
	if err := ctx.BindJSON(&input); err != nil {
		fmt.Println("Error: Invalid JSON type", err)
	}

	now := time.Now()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO products (name, price, description, product_size, stock, isFlashSale, isfavorite_product, temperature, category_product_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, input.Name, input.Price, input.Description, input.Productsize, input.Stock, input.Isflashsale, input.IsFavorite_product, input.Tempelatur, input.Category_productid, now, now)

	if err != nil {
		fmt.Println("Error inserting product:", err)
	}

	input.Created_at = now
	input.Updated_at = now
	return input
}

func GetById(ctx *gin.Context, pool *pgxpool.Pool) (models.Product, error) {
	id := ctx.Param("id")
	var product models.Product

	err := pool.QueryRow(context.Background(), `
		SELECT id, name, price, description, product_size, stock, isFlashSale, isFavorite_product, temperature, category_product_id, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Productsize, &product.Stock, &product.Isflashsale, &product.IsFavorite_product, &product.Tempelatur, &product.Category_productid, &product.Created_at, &product.Updated_at)

	return product, err
}

func Edit(pool *pgxpool.Pool, ctx *gin.Context) (models.Product, error) {
	id := ctx.Param("id")

	oldProduct, err := GetById(ctx, pool)

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
	if newProduct.IsFavorite_product == nil {
		newProduct.IsFavorite_product = oldProduct.IsFavorite_product
	}
	if newProduct.Tempelatur == "" {
		newProduct.Tempelatur = oldProduct.Tempelatur
	}
	if newProduct.Category_productid == 0 {
		newProduct.Category_productid = oldProduct.Category_productid
	}

	_, err = pool.Exec(context.Background(),
		`UPDATE products 
     SET name=$1, price=$2, description=$3, product_size=$4, stock=$5, 
         isflashsale=$6, isfavorite_product=$7, temperature=$8, category_product_id=$9, updated_at=NOW() 
     WHERE id = $10`,
		newProduct.Name, newProduct.Price, newProduct.Description, newProduct.Productsize,
		newProduct.Stock, *newProduct.Isflashsale, *newProduct.IsFavorite_product,
		newProduct.Tempelatur, newProduct.Category_productid, id)

	return newProduct, err
}

func Delete(pool *pgxpool.Pool, ctx *gin.Context) error {
	id := ctx.Param("id")
	_, err := pool.Exec(context.Background(), "DELETE FROM products WHERE id = $1", id)

	return err
}

func CreateImageProduct(pool *pgxpool.Pool, ctx *gin.Context, productId int, files []*multipart.FileHeader) ([]models.ImageProduct, error) {
	var inputImage []models.ImageProduct
	now := time.Now()
	maxSize := int64(5 * 1024 * 1024)
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	for _, file := range files {
		if file.Size > maxSize {
			fmt.Printf("file %s melebihi ukuran maksimum 5MB", file.Filename)
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedTypes[ext] {
			fmt.Printf("file %s bukan tipe gambar yang diizinkan (hanya jpg, jpeg, png, webp)", file.Filename)
		}

		err := os.MkdirAll("imagesProduct", os.ModePerm)
		if err != nil {
			fmt.Println("Error : Failed to create folder", err)
		}

		filePath := fmt.Sprintf("imagesProduct/%s", file.Filename)

		err = saveUploadedFile(file, filePath)
		if err != nil {
			fmt.Println("Error :", err)
		}

		_, err = pool.Exec(context.Background(), "INSERT INTO imageproduct (product_id, image, created_at, updated_at) VALUES ($1, $2, $3, $4)", productId, filePath, now, now)

		if err != nil {
			fmt.Println("Error: Failed to create image prodct", err)
		}

		inputImage = append(inputImage, models.ImageProduct{
			Productid:  productId,
			Image:      filePath,
			Created_at: now,
			Updated_at: now,
		})
	}
	return inputImage, nil
}

func GetAllImageProduct(pool *pgxpool.Pool) ([]models.ImageProduct, error) {
	var images []models.ImageProduct

	rows, err := pool.Query(context.Background(), "SELECT id, product_id, image, created_at, updated_at  FROM image_products ORDER BY id ASC")

	if err != nil {
		fmt.Println("Error : Failed to get all image product", err)
	}

	for rows.Next() {
		var img models.ImageProduct
		err := rows.Scan(&img.Id, &img.Productid, &img.Image, &img.Created_at, &img.Updated_at)
		if err != nil {
			fmt.Println("Error scanning image product:", err)
		}
		images = append(images, img)
	}

	return images, nil
}

func DeleteImageProduct(pool *pgxpool.Pool, id int) error {
	var imagePath string

	err := pool.QueryRow(context.Background(), "SELECT image FROM image_products WHERE id = $1", id).Scan(&imagePath)
	if err != nil {
		fmt.Println("image not found:", err)
	}

	err = os.Remove(imagePath)
	if err != nil {
		fmt.Println("Error: Failed to delete image:", err)
	}

	_, err = pool.Exec(context.Background(), "DELETE FROM image_products WHERE id = $1", id)
	if err != nil {
		fmt.Println("failed to delete image product:", err)
	}

	return nil
}

func saveUploadedFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func GetProductFavorite(pool *pgxpool.Pool, limit int) ([]models.Product, int, error) {
	var products []models.Product

	if limit < 4 {
		limit = 4
	}

	var total int
	err := pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM products WHERE isFavorite_product = true").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting favorite products: %v", err)
	}

	rows, err := pool.Query(context.Background(), `
		SELECT 
			p.id,
			p.name,
			p.description,
			p.price,
			COALESCE(ip.image, '') AS image
		FROM products p
		LEFT JOIN image_products ip ON ip.product_id = p.id
		WHERE p.isFavorite_product = true
		GROUP BY p.id, ip.image
		ORDER BY p.id
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching favorite products: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Image)
		if err != nil {
			fmt.Println("Error scanning product:", err)
			continue
		}
		products = append(products, p)
	}

	return products, total, nil
}

func FilterProducts(pool *pgxpool.Pool, name, categoryStr, sortBy, priceMin, priceMax string, page, limit int) ([]models.Product, int, error) {
	var products []models.Product
	offset := (page - 1) * limit

	query := `
	SELECT 
			p.id, p.name, p.price, p.description, 
			p.product_size, p.stock, p.isFlashSale, p.isFavorite_product, 
			p.temperature, p.category_product_id, COALESCE(ip.image, '') AS image
		FROM products p
		LEFT JOIN image_products ip ON ip.product_id = p.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if name != "" {
		query += fmt.Sprintf("AND LOWER(p.name) LIKE LOWER($%d)", argIndex)
		args = append(args, "%"+name+"%")
		argIndex++
	}

	if categoryStr != "" {
		categorys := strings.Split(categoryStr, ",")
		query += fmt.Sprintf(" AND p.category_product_id = ANY($%d)", argIndex)
		args = append(args, categorys)
		argIndex++
	}

	if priceMin != "" {
		query += fmt.Sprintf(" AND p.price >= $%d", argIndex)
		args = append(args, priceMin)
		argIndex++
	}

	if priceMax != "" {
		query += fmt.Sprintf(" AND p.price <= $%d", argIndex)
		args = append(args, priceMax)
		argIndex++
	}

	switch sortBy {
	case "price_asc":
		query += " ORDER BY p.price ASC"
	case "price_desc":
		query += " ORDER BY p.price DESC"
	case "newest":
		query += " ORDER BY p.created_at DESC"
	default:
		query += " ORDER BY p.id ASC"
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Description, &p.Productsize, &p.Stock, &p.Isflashsale, &p.IsFavorite_product, &p.Tempelatur, &p.Category_productid, &p.Image)
		if err != nil {
			fmt.Println("Error scanning product:", err)
			continue
		}
		products = append(products, p)
	}

	countQuery := "SELECT COUNT(*) FROM products p WHERE 1=1"
	countArgs := []interface{}{}
	argIdx := 1

	if name != "" {
		countQuery += fmt.Sprintf(" AND LOWER(p.name) LIKE LOWER($%d)", argIdx)
		countArgs = append(countArgs, "%"+name+"%")
		argIdx++
	}
	if categoryStr != "" {
		categories := strings.Split(categoryStr, ",")
		countQuery += fmt.Sprintf(" AND p.category_product_id = ANY($%d)", argIdx)
		countArgs = append(countArgs, categories)
		argIdx++
	}
	if priceMin != "" {
		countQuery += fmt.Sprintf(" AND p.price >= $%d", argIdx)
		countArgs = append(countArgs, priceMin)
		argIdx++
	}
	if priceMax != "" {
		countQuery += fmt.Sprintf(" AND p.price <= $%d", argIdx)
		countArgs = append(countArgs, priceMax)
		argIdx++
	}

	var total int
	err = pool.QueryRow(context.Background(), countQuery, countArgs...).Scan(&total)
	if err != nil {
		fmt.Println("Error :", err)
	}

	return products, total, nil
}
