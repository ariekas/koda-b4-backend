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
	var products []models.Product
	limit := 50
	offset := (page - 1) * limit

	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		fmt.Println("Error counting products:", err)
	}

	rows, err := pool.Query(context.Background(), `
	SELECT p.id, p.name, p.price, p.description, p.stock, 
	       p.is_flashsale, p.is_favorite_product, p.category_products_id, 
	       COALESCE(ip.image, '') AS image,
	       p.created_at, p.updated_at
	FROM products p
	LEFT JOIN product_images ip ON ip.products_id = p.id
	ORDER BY p.id
	OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		fmt.Println("Error: Failed get data product", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Description, &p.Stock,
			&p.IsFlashSale, &p.IsFavoriteProduct, &p.CategoryProductId, &p.Image,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			fmt.Println("Error scanning product:", err)
			continue
		}
		products = append(products, p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	links := map[string]string{}
	if page > 1 {
		links["prev"] = fmt.Sprintf("/products?page=%d", page-1)
	} else {
		links["prev"] = "null"
	}
	if page < totalPages {
		links["next"] = fmt.Sprintf("/products?page=%d", page+1)
	} else {
		links["next"] = "null"
	}

	return PaginationResponse{
		Data:       products,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}, nil
}

func CreateProduct(pool *pgxpool.Pool, input models.ProductInput) (models.Product, error) {
	now := time.Now()

	var discountsID interface{} = nil
	priceDiscount := 0.0

	_, err := pool.Exec(context.Background(), `
		INSERT INTO products 
		(discounts_id, name, price, price_discounts, description, stock, is_flashsale, 
		 is_favorite_product, category_products_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`,
		discountsID,             
		input.Name,                
		input.Price,                
		priceDiscount,              
		input.Description,         
		input.Stock,              
		input.IsFlashSale,         
		false,                      
		input.CategoryProductId,   
		now, now,
	)

	if err != nil {
		fmt.Println("failed to insert product:", err)
		return models.Product{}, err
	}

	var productId int
	err = pool.QueryRow(context.Background(),
		"SELECT id FROM products WHERE name=$1 ORDER BY id DESC LIMIT 1",
		input.Name,
	).Scan(&productId)

	if err != nil {
		fmt.Println("Error fetching product ID:", err)
		return models.Product{}, err
	}

	product := models.Product{
		Id:                productId,
		Name:              input.Name,
		Price:             input.Price,
		Description:       input.Description,
		Stock:             input.Stock,
		IsFlashSale:       *input.IsFlashSale,
		IsFavoriteProduct: false,
		CategoryProductId: input.CategoryProductId,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	return product, nil
}

func GetProductByID(pool *pgxpool.Pool, id int) (models.Product, error) {
	var p models.Product
	err := pool.QueryRow(context.Background(), `
	SELECT id, name, price, description, stock, is_flashsale, is_favorite_product, category_products_id, created_at, updated_at
	FROM products WHERE id=$1
	`, id).Scan(&p.Id, &p.Name, &p.Price, &p.Description, &p.Stock, &p.IsFlashSale, &p.IsFavoriteProduct, &p.CategoryProductId, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		fmt.Println("Error : Failed to get product by id", err)
	}
	return p, nil
}

func EditProduct(pool *pgxpool.Pool, id int, input models.ProductInput) (models.Product, error) {
	old, err := GetProductByID(pool, id)
	if err != nil {
		return models.Product{}, fmt.Errorf("product not found")
	}

	if input.Name == "" {
		input.Name = old.Name
	}
	if input.Price == 0 {
		input.Price = old.Price
	}
	if input.Description == "" {
		input.Description = old.Description
	}
	if input.Stock == 0 {
		input.Stock = old.Stock
	}
	if input.IsFlashSale == nil {
		input.IsFlashSale = &old.IsFlashSale
	}
	if input.CategoryProductId == 0 {
		input.CategoryProductId = old.CategoryProductId
	}

	_, err = pool.Exec(context.Background(), `
	UPDATE products 
	SET name=$1, price=$2, description=$3, stock=$4, is_flashsale=$5, category_products_id=$6, updated_at=NOW()
	WHERE id=$7
	`, input.Name, input.Price, input.Description, input.Stock, *input.IsFlashSale, input.CategoryProductId, id)
	if err != nil {
		return models.Product{}, err
	}

	return GetProductByID(pool, id)
}

func DeleteProduct(pool *pgxpool.Pool, id int) error {
	_, err := pool.Exec(context.Background(), "DELETE FROM products WHERE id=$1", id)
	return err
}

func CreateImageProduct(pool *pgxpool.Pool, productId int, files []*multipart.FileHeader) ([]models.ImageProduct, error) {
	var images []models.ImageProduct
	now := time.Now()
	maxSize := int64(5 * 1024 * 1024)
	allowedTypes := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}

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

		err = SaveUploadedFile(file, filePath)
		if err != nil {
			fmt.Println("Error :", err)
		}

		_, err = pool.Exec(context.Background(), `
		INSERT INTO product_images (products_id, image, created_at, updated_at)
		VALUES ($1,$2,$3,$4)
		`, productId, filePath, now, now)
		if err != nil {
			fmt.Println("Error: Failed to create image prodct", err)
		}

		images = append(images, models.ImageProduct{
			ProductId: productId,
			Image:     filePath,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return images, nil
}

func SaveUploadedFile(file *multipart.FileHeader, path string) error {
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

func GetAllImageProducts(pool *pgxpool.Pool) ([]models.ImageProduct, error) {
	var images []models.ImageProduct
	rows, err := pool.Query(context.Background(), "SELECT id, products_id, image, created_at, updated_at FROM product_images")
	if err != nil {
		return images, err
	}
	defer rows.Close()
	for rows.Next() {
		var img models.ImageProduct
		rows.Scan(&img.Id, &img.ProductId, &img.Image, &img.CreatedAt, &img.UpdatedAt)
		images = append(images, img)
	}
	return images, nil
}


func DeleteImageProduct(pool *pgxpool.Pool, id int) error {
	var path string
	err := pool.QueryRow(context.Background(), "SELECT image FROM product_images WHERE id=$1", id).Scan(&path)
	if err != nil {
		return err
	}
	os.Remove(path)
	_, err = pool.Exec(context.Background(), "DELETE FROM product_images WHERE id=$1", id)
	return err
}

func GetProductFavorite(pool *pgxpool.Pool, limit int) ([]models.Product, int, error) {
	if limit < 1 {
		limit = 4
	}
	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM products WHERE is_favorite_product=true").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting favorite products: %v", err)
	}

	rows, err := pool.Query(context.Background(), `
	SELECT p.id, p.name, p.price, p.description, p.stock, p.is_flashsale, p.is_favorite_product, p.category_products_id, COALESCE(ip.image,'') as image, p.created_at, p.updated_at
	FROM products p
	LEFT JOIN product_images ip ON ip.products_id=p.id
	WHERE p.is_favorite_product=true
	ORDER BY p.id ASC
	LIMIT $1
	`, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching favorite products: %v", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.Id, &p.Name, &p.Price, &p.Description, &p.Stock, &p.IsFlashSale, &p.IsFavoriteProduct, &p.CategoryProductId, &p.Image, &p.CreatedAt, &p.UpdatedAt)
		products = append(products, p)
	}
	return products, total, nil
}

func FilterProducts(pool *pgxpool.Pool, name, categoryStr, sortBy string, priceMin, priceMax float64, page, limit int) ([]models.Product, int, error) {
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := `
	SELECT p.id, p.name, p.price, p.description, p.stock, p.is_flashsale, p.is_favorite_product, p.category_products_id, COALESCE(ip.image,'') AS image, p.created_at, p.updated_at
	FROM products p
	LEFT JOIN product_images ip ON ip.products_id=p.id
	WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1
	if name != "" {
		query += fmt.Sprintf(" AND LOWER(p.name) LIKE LOWER($%d)", argIdx)
		args = append(args, "%"+name+"%")
		argIdx++
	}

	if categoryStr != "" {
		categories := strings.Split(categoryStr, ",")
		query += fmt.Sprintf(" AND p.category_products_id = ANY($%d)", argIdx)
		args = append(args, categories)
		argIdx++
	}


	if priceMin > 0 {
		query += fmt.Sprintf(" AND p.price >= $%d", argIdx)
		args = append(args, priceMin)
		argIdx++
	}
	if priceMax > 0 {
		query += fmt.Sprintf(" AND p.price <= $%d", argIdx)
		args = append(args, priceMax)
		argIdx++
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

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch filtered products: %v", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Description, &p.Stock,
			&p.IsFlashSale, &p.IsFavoriteProduct, &p.CategoryProductId, &p.Image, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			fmt.Println("Error scanning product:", err)
			continue
		}
		products = append(products, p)
	}

	countQuery := "SELECT COUNT(*) FROM products p WHERE 1=1"
	countArgs := []interface{}{}
	argIdx = 1

	if name != "" {
		countQuery += fmt.Sprintf(" AND LOWER(p.name) LIKE LOWER($%d)", argIdx)
		countArgs = append(countArgs, "%"+name+"%")
		argIdx++
	}
	if categoryStr != "" {
		categories := strings.Split(categoryStr, ",")
		countQuery += fmt.Sprintf(" AND p.category_products_id = ANY($%d)", argIdx)
		countArgs = append(countArgs, categories)
		argIdx++
	}
	if priceMin > 0 {
		countQuery += fmt.Sprintf(" AND p.price >= $%d", argIdx)
		countArgs = append(countArgs, priceMin)
		argIdx++
	}
	if priceMax > 0 {
		countQuery += fmt.Sprintf(" AND p.price <= $%d", argIdx)
		countArgs = append(countArgs, priceMax)
		argIdx++
	}

	var total int
	err = pool.QueryRow(context.Background(), countQuery, countArgs...).Scan(&total)
	if err != nil {
		fmt.Println("Error counting filtered products:", err)
	}

	return products, total, nil
}


func DetailProduct(pool *pgxpool.Pool, id int) (models.ProductDetail, error) {
	var detail models.ProductDetail
	product, err := GetProductByID(pool, id)
	if err != nil {
		return detail, err
	}

	detail.Product = product

	rowsImg, _ := pool.Query(context.Background(), "SELECT id, products_id, image, created_at, updated_at FROM product_images WHERE products_id=$1", id)
	var images []models.ImageProduct
	for rowsImg.Next() {
		var img models.ImageProduct
		rowsImg.Scan(&img.Id, &img.ProductId, &img.Image, &img.CreatedAt, &img.UpdatedAt)
		images = append(images, img)
	}
	detail.Images = images

	rowsSize, _ := pool.Query(context.Background(), "SELECT id, name, product_id, created_at, updated_at FROM size_products WHERE product_id=$1", id)
	var sizes []models.SizeProduct
	for rowsSize.Next() {
		var s models.SizeProduct
		rowsSize.Scan(&s.Id, &s.Name, &s.ProductId, &s.CreatedAt, &s.UpdatedAt)
		sizes = append(sizes, s)
	}
	detail.Sizes = sizes

	rowsVariant, _ := pool.Query(context.Background(), "SELECT id, name, product_id, created_at, updated_at FROM variant_products WHERE product_id=$1", id)
	var variants []models.VariantProduct
	for rowsVariant.Next() {
		var v models.VariantProduct
		rowsVariant.Scan(&v.Id, &v.Name, &v.ProductId, &v.CreatedAt, &v.UpdatedAt)
		variants = append(variants, v)
	}
	detail.Variants = variants

	return detail, nil
}