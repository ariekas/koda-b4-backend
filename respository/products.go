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
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM product").Scan(&total)
	if err != nil {
		fmt.Println("Error counting products:", err)
	}

	rows, err := pool.Query(context.Background(), `
	SELECT id, name, price, description, productsize, stock, isflashsale, tempelatur, category_productid, created_at, updated_at 
	FROM product 
	ORDER BY id 
	OFFSET $1 LIMIT $2
`, offset, limit)
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

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	links := make(map[string]string)

	if page > 1 {
		links["prev"] = fmt.Sprintf("/products?page=%d", page-1)
	} else {
		links["prev"] = "nul"
	}

	if page < totalPages {
		links["next"] = fmt.Sprintf("/products?page=%d", page+1)
	}

	response := PaginationResponse{
		Data:       dataProduct,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}

	return response, nil
}

func Create(ctx *gin.Context, pool *pgxpool.Pool) models.Product {
	var input models.Product

	err := ctx.BindJSON(&input)

	if err != nil {
		fmt.Println("Error: Invalid type much json")
	}

	now := time.Now()

	_, err = pool.Exec(context.Background(), "INSERT INTO product (name, price, description, productsize, stock, isflashsale, tempelatur, category_productid, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", input.Name, input.Price, input.Description, input.Productsize, input.Stock, input.Isflashsale, input.Tempelatur, input.Category_productid, now, now)

	if err != nil {
		fmt.Println("Error insert product:", err)
	}

	input.Created_at = now
	input.Updated_at = now

	return input
}

func GetById(ctx *gin.Context, pool *pgxpool.Pool) (models.Product, error) {
	id := ctx.Param("id")

	var product models.Product

	err := pool.QueryRow(context.Background(), `
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
	if newProduct.Tempelatur == "" {
		newProduct.Tempelatur = oldProduct.Tempelatur
	}
	if newProduct.Category_productid == 0 {
		newProduct.Category_productid = oldProduct.Category_productid
	}

	_, err = pool.Exec(context.Background(), "UPDATE product SET name=$1, price=$2, description=$3, productsize=$4, stock=$5, isflashsale=$6, tempelatur=$7, category_productid=$8, updated_at=NOW() WHERE id = $9", newProduct.Name, newProduct.Price, newProduct.Description, newProduct.Productsize,
		newProduct.Stock, *newProduct.Isflashsale, newProduct.Tempelatur, newProduct.Category_productid, id)

	return newProduct, err
}

func Delete(pool *pgxpool.Pool, ctx *gin.Context) error {
	id := ctx.Param("id")
	_, err := pool.Exec(context.Background(), "DELETE FROM product WHERE id = $1", id)

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

		_, err = pool.Exec(context.Background(), "INSERT INTO imageproduct (productid, image, created_at, updated_at) VALUES ($1, $2, $3, $4)", productId, filePath, now, now)

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

	rows, err := pool.Query(context.Background(), "SELECT id, productid, image, created_at, updated_at  FROM imageproduct ORDER BY id ASC")

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

func DeleteImageProduct(pool *pgxpool.Pool, id int) error{
	var imagePath string

	err := pool.QueryRow(context.Background(), "SELECT image FROM imageproduct WHERE id = $1", id).Scan(&imagePath)
	if err != nil {
		fmt.Println("image not found:", err)
	}

	err = os.Remove(imagePath)
	if err != nil {
		fmt.Println("Error: Failed to delete image:", err)
	}

	_, err = pool.Exec(context.Background(), "DELETE FROM imageproduct WHERE id = $1", id)
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
