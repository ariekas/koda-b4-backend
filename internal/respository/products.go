package respository

import (
	"back-end-coffeShop/internal/config"
	"back-end-coffeShop/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetProducts(pool *pgxpool.Pool, page int) (models.PaginationResponse, error) {
	var products []models.Product
	limit := 50
	offset := (page - 1) * limit

	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		fmt.Println("Error counting products:", err)
		return models.PaginationResponse{}, err
	}

	rows, err := pool.Query(context.Background(), `
	SELECT 
	    p.id, 
	    p.name, 
	    p.price, 
	    p.description, 
	    p.stock, 
	    p.is_flashsale, 
	    p.is_favorite_product, 
	    p.category_products_id,
	    COALESCE(
	        JSON_AGG(
                JSON_BUILD_OBJECT('id', ip.id, 'image', ip.image)
	        ) FILTER (WHERE ip.id IS NOT NULL), '[]'
	    ) AS images,
	    p.created_at, 
	    p.updated_at
	FROM products p
	LEFT JOIN product_images ip ON ip.products_id = p.id
	GROUP BY p.id
	OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		fmt.Println("Error: Failed get data product", err)
		return models.PaginationResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		var imagesJSON []byte

		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.Description,
			&p.Stock,
			&p.IsFlashSale,
			&p.IsFavoriteProduct,
			&p.CategoryProductId,
			&imagesJSON,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			fmt.Println("Error scanning product:", err)
			continue
		}

		if err := json.Unmarshal(imagesJSON, &p.Images); err != nil {
			fmt.Println("Error unmarshalling images:", err)
			p.Images = []models.ImageProduct{}
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

	return models.PaginationResponse{
		Data:       products,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}, nil
}
func CreateProduct(pool *pgxpool.Pool, input models.ProductCreateInput) (*models.Product, error) {
	ctx := context.Background()
	now := time.Now()

	priceDiscount := 0.0

	if input.DiscountsId != nil {
		var discount float64
		err := pool.QueryRow(ctx,
			"SELECT discount_percentage FROM discounts WHERE id = $1",
			input.DiscountsId,
		).Scan(&discount)

		if err == nil {
			priceDiscount = input.Price - (input.Price * (discount / 100))
		}
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var newID int

	err = tx.QueryRow(ctx, `
        INSERT INTO products (
            discounts_id,
            name,
            price,
            price_discount,
            description,
            stock,
            is_flashsale,
            is_favorite_product,
            category_products_id,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id
    `,
		input.DiscountsId,
		input.Name,
		input.Price,
		priceDiscount, 
		input.Description,
		input.Stock,
		input.IsFlashSale,
		input.IsFavoriteProduct,
		input.CategoryProductId,
		now,
		now,
	).Scan(&newID)

	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	for _, sizeId := range input.SizeIds {
		_, err = tx.Exec(ctx, `
            INSERT INTO product_sizes (products_id, size_products_id, created_at, updated_at)
            VALUES ($1, $2, $3, $4)
        `, newID, sizeId, now, now)
		if err != nil {
			return nil, fmt.Errorf("failed to insert product size: %w", err)
		}
	}

	for _, variantId := range input.VariantIds {
		_, err = tx.Exec(ctx, `
            INSERT INTO product_variants (products_id, variant_products_id, created_at, updated_at)
            VALUES ($1, $2, $3, $4)
        `, newID, variantId, now, now)
		if err != nil {
			return nil, fmt.Errorf("failed to insert product variant: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	product := models.Product{
		Id:                newID,
		Name:              input.Name,
		Price:             input.Price,
		PriceDiscount:     priceDiscount,
		Description:       input.Description,
		Stock:             input.Stock,
		IsFlashSale:       input.IsFlashSale,
		IsFavoriteProduct: input.IsFavoriteProduct,
		CategoryProductId: input.CategoryProductId,
		DiscountsId:       input.DiscountsId,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	return &product, nil
}

func GetProductByID(pool *pgxpool.Pool, id int) (models.Product, error) {
	var prodct models.Product

	err := pool.QueryRow(context.Background(), `
		SELECT 
			id,
			name,
			price,
			price_discount,
			discounts_id,
			description,
			stock,
			is_flashsale,
			is_favorite_product,
			category_products_id,
			created_at,
			updated_at
		FROM products 
		WHERE id = $1
	`, id).Scan(
		&prodct.Id,
		&prodct.Name,
		&prodct.Price,
		&prodct.PriceDiscount,
		&prodct.DiscountsId,
		&prodct.Description,
		&prodct.Stock,
		&prodct.IsFlashSale,
		&prodct.IsFavoriteProduct,
		&prodct.CategoryProductId,
		&prodct.CreatedAt,
		&prodct.UpdatedAt,
	)

	if err != nil {
		return models.Product{}, fmt.Errorf("error: Failed to get product id, %w", err)
	}

	return prodct, nil
}

func GetDiscountByID(pool *pgxpool.Pool, id int) (models.Discount, error) {
	row := pool.QueryRow(context.Background(), `
		SELECT id, diskon FROM discounts WHERE id=$1
	`, id)

	var d models.Discount
	err := row.Scan(&d.Id, &d.Diskon)
	return d, err
}

func EditProduct(pool *pgxpool.Pool, id int, input models.ProductUpdateInput) error {
	oldProduct, err := GetProductByID(pool, id)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	name := oldProduct.Name
	if input.Name != nil {
		name = *input.Name
	}

	price := oldProduct.Price
	if input.Price != nil {
		price = *input.Price
	}

	description := oldProduct.Description
	if input.Description != nil {
		description = *input.Description
	}

	stock := oldProduct.Stock
	if input.Stock != nil {
		stock = *input.Stock
	}

	isFlashSale := oldProduct.IsFlashSale
	if input.IsFlashSale != nil {
		isFlashSale = *input.IsFlashSale
	}

	isFavoriteProduct := oldProduct.IsFavoriteProduct
	if input.IsFavoriteProduct != nil {
		isFavoriteProduct = *input.IsFavoriteProduct
	}

	categoryProductId := oldProduct.CategoryProductId
	if input.CategoryProductId != nil {
		categoryProductId = *input.CategoryProductId
	}

	var discountsId *int

	if input.DiscountsId != nil {
		discountsId = input.DiscountsId
	} else {
		discountsId = oldProduct.DiscountsId
	}

	priceDiscount := price
	if discountsId != nil && *discountsId > 0 {
		discount, err := GetDiscountByID(pool, *discountsId)
		if err == nil {
			discountAmount := price * float64(discount.Diskon) / 100
			priceDiscount = price - discountAmount
		}
	}

	query := `
		UPDATE products 
		SET 
			name=$1,
			price=$2,
			price_discount=$3,
			discounts_id=$4,
			description=$5,
			stock=$6,
			is_flashsale=$7,
			is_favorite_product=$8,
			category_products_id=$9,
			updated_at=NOW()
		WHERE id=$10
	`

	_, err = pool.Exec(context.Background(), query,
		name,
		price,
		priceDiscount,
		discountsId,
		description,
		stock,
		isFlashSale,
		isFavoriteProduct,
		categoryProductId,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func DeleteProduct(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()

	_, err := pool.Exec(ctx,
		"DELETE FROM product_sizes WHERE products_id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete product sizes: %w", err)
	}

	_, err = pool.Exec(ctx,
		"DELETE FROM product_variants WHERE products_id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete product variants: %w", err)
	}

	_, err = pool.Exec(ctx,
		"DELETE FROM product_images WHERE products_id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete product images: %w", err)
	}

	_, err = pool.Exec(ctx,
		"DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func CreateImageProduct(pool *pgxpool.Pool, ctx *gin.Context, productId int, files []*multipart.FileHeader) ([]models.ImageProduct, error) {
	var images []models.ImageProduct
	now := time.Now()

	var exists bool
	err := pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)",
		productId,
	).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("failed checking product: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("product not found")
	}

	maxSize := int64(5 * 1024 * 1024)
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	localPath := "uploads/imageProducts"

	useCloud := os.Getenv("CLOUDINARY_API_KEY") != "" &&
		os.Getenv("CLOUDINARY_API_SECRET") != "" &&
		os.Getenv("CLOUDINARY_NAME") != ""

	if !useCloud {
		os.MkdirAll(localPath, 0755)
	}

	for _, file := range files {

		if file.Size > maxSize {
			return nil, fmt.Errorf("file %s melebihi 5MB", file.Filename)
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedTypes[ext] {
			return nil, fmt.Errorf("file %s bukan tipe gambar valid", file.Filename)
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		var finalURL string

		if useCloud {
			imageURL, err := config.UploaderFile(file, "imageProducts", fileName)
			if err != nil {
				return nil, fmt.Errorf("cloud upload failed: %v", err)
			}
			finalURL = imageURL

		} else {
			savePath := filepath.Join(localPath, fileName)

			if err := ctx.SaveUploadedFile(file, savePath); err != nil {
				return nil, fmt.Errorf("failed upload local: %v", err)
			}

			finalURL = "uploads/imageProducts/" + fileName
		}

		var newID int
		err = pool.QueryRow(context.Background(), `
			INSERT INTO product_images (products_id, image, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, productId, finalURL, now, now).Scan(&newID)

		if err != nil {
			return nil, fmt.Errorf("failed to save image db: %v", err)
		}

		images = append(images, models.ImageProduct{
			Id:        newID,
			ProductId: productId,
			Image:     finalURL,
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
		return fmt.Errorf("image not found")
	}

	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	_, err = pool.Exec(context.Background(), "DELETE FROM product_images WHERE id=$1", id)
	return err
}

func GetProductFavorite(pool *pgxpool.Pool, limit int) ([]models.Product, int, error) {
	if limit < 1 {
		limit = 4
	}

	var total int
	err := pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM products WHERE is_favorite_product=true").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting favorite products: %v", err)
	}

	rows, err := pool.Query(context.Background(), `
		SELECT 
			p.id,
			p.name,
			p.price,
			p.description,
			p.stock,
			p.is_flashsale,
			p.is_favorite_product,
			p.category_products_id,
			COALESCE((
				SELECT pi.image 
				FROM product_images pi 
				WHERE pi.products_id = p.id 
				ORDER BY pi.id ASC 
				LIMIT 1
			), '') AS image,
			p.created_at,
			p.updated_at
		FROM products p
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
		var image string

		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.Description,
			&p.Stock,
			&p.IsFlashSale,
			&p.IsFavoriteProduct,
			&p.CategoryProductId,
			&image,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			fmt.Println("scan error:", err)
			continue
		}

		p.Images = []models.ImageProduct{
			{Image: image},
		}

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
	SELECT 
		p.id,
		p.name,
		p.price,
		p.description,
		p.stock,
		p.is_flashsale,
		p.is_favorite_product,
		p.category_products_id,
		COALESCE((
			SELECT pi.image 
			FROM product_images pi 
			WHERE pi.products_id = p.id 
			ORDER BY pi.id ASC 
			LIMIT 1
		), '') AS image,
		p.created_at,
		p.updated_at
	FROM products p
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
		split := strings.Split(categoryStr, ",")
		categories := make([]int, 0)
		for _, c := range split {
			if v, err := strconv.Atoi(c); err == nil {
				categories = append(categories, v)
			}
		}

		if len(categories) > 0 {
			query += fmt.Sprintf(" AND p.category_products_id = ANY($%d)", argIdx)
			args = append(args, categories)
			argIdx++
		}
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
		var image string

		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.Description,
			&p.Stock,
			&p.IsFlashSale,
			&p.IsFavoriteProduct,
			&p.CategoryProductId,
			&image,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			fmt.Println("Error scanning product:", err)
			continue
		}

		p.Images = []models.ImageProduct{
			{Image: image},
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
		split := strings.Split(categoryStr, ",")
		categories := make([]int, 0)
		for _, c := range split {
			if v, err := strconv.Atoi(c); err == nil {
				categories = append(categories, v)
			}
		}

		if len(categories) > 0 {
			countQuery += fmt.Sprintf(" AND p.category_products_id = ANY($%d)", argIdx)
			countArgs = append(countArgs, categories)
			argIdx++
		}
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
		fmt.Println("Error counting products:", err)
	}

	return products, total, nil
}

func DetailProduct(pool *pgxpool.Pool, id int) (models.Product, error) {
	var detail models.Product

	product, err := GetProductByID(pool, id)
	if err != nil {
		return detail, err
	}
	detail = product

	rowsImg, err := pool.Query(context.Background(),
		`
		SELECT id, products_id, image, created_at, updated_at 
		FROM product_images 
		WHERE products_id = $1
		`, id)
	if err != nil {
		return detail, fmt.Errorf("error fetching product images: %v", err)
	}
	defer rowsImg.Close()

	var images []models.ImageProduct
	for rowsImg.Next() {
		var img models.ImageProduct
		if err := rowsImg.Scan(&img.Id, &img.ProductId, &img.Image, &img.CreatedAt, &img.UpdatedAt); err != nil {
			return detail, fmt.Errorf("error scanning product image: %v", err)
		}
		images = append(images, img)
	}
	detail.Images = images

	rowsSize, err := pool.Query(context.Background(),
		`
		SELECT id, name, additional_costs, created_at, updated_at 
		FROM size_products
		`)
	if err != nil {
		return detail, fmt.Errorf("error fetching product sizes: %v", err)
	}
	defer rowsSize.Close()

	var sizes []models.SizeProduct
	for rowsSize.Next() {
		var s models.SizeProduct
		if err := rowsSize.Scan(&s.Id, &s.Name, &s.AdditionalCosts, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return detail, fmt.Errorf("error scanning product size: %v", err)
		}
		sizes = append(sizes, s)
	}
	detail.Sizes = sizes

	rowsVariant, err := pool.Query(context.Background(),
		`
		SELECT id, name, additional_costs, created_at, updated_at 
		FROM variant_products
		`)
	if err != nil {
		return detail, fmt.Errorf("error fetching product variants: %v", err)
	}
	defer rowsVariant.Close()

	var variants []models.VariantProduct
	for rowsVariant.Next() {
		var v models.VariantProduct
		if err := rowsVariant.Scan(&v.Id, &v.Name, &v.AdditionalCosts, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return detail, fmt.Errorf("error scanning product variant: %v", err)
		}
		variants = append(variants, v)
	}
	detail.Variants = variants

	return detail, nil
}
