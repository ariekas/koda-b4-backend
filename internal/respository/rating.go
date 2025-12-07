package respository

import (
	"back-end-coffeShop/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateRating(pool *pgxpool.Pool, rating *models.Rating) error {

	now := time.Now()
	err := pool.QueryRow(context.Background(),
		`
		INSERT INTO ratings (users_id, products_id, transactions_id, rating, review, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		rating.UsersID,
		rating.ProductsID,
		rating.TransactionsID,
		rating.Rating,
		rating.Review,
		now,
		now,
	).Scan(&rating.ID)

	if err != nil {
		return err
	}

	rating.CreatedAt = now
	rating.UpdatedAt = now
	return nil
}

func DetailRating(pool *pgxpool.Pool, id int) (*models.Rating, error) {

	rating := &models.Rating{}
	err := pool.QueryRow(context.Background(), `
	SELECT id, users_id, products_id, transactions_id, rating, review, created_at, updated_at
	FROM ratings
	WHERE id = $1
`, id).Scan(
		&rating.ID,
		&rating.UsersID,
		&rating.ProductsID,
		&rating.TransactionsID,
		&rating.Rating,
		&rating.Review,
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("rating not found")
	}

	return rating, err
}

func GetRatingProduct(pool *pgxpool.Pool, productID int) ([]models.Rating, error) {

	rows, err := pool.Query(context.Background(), `
	SELECT id, users_id, products_id, transactions_id, rating, review, created_at, updated_at
	FROM ratings
	WHERE products_id = $1
	ORDER BY created_at DESC
`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []models.Rating
	for rows.Next() {
		var rating models.Rating
		err := rows.Scan(
			&rating.ID,
			&rating.UsersID,
			&rating.ProductsID,
			&rating.TransactionsID,
			&rating.Rating,
			&rating.Review,
			&rating.CreatedAt,
			&rating.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func GetRatingUser(pool *pgxpool.Pool, userID int) ([]models.Rating, error) {

	rows, err := pool.Query(context.Background(), `
	SELECT id, users_id, products_id, transactions_id, rating, review, created_at, updated_at
	FROM ratings
	WHERE users_id = $1
	ORDER BY created_at DESC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []models.Rating
	for rows.Next() {
		var rating models.Rating
		err := rows.Scan(
			&rating.ID,
			&rating.UsersID,
			&rating.ProductsID,
			&rating.TransactionsID,
			&rating.Rating,
			&rating.Review,
			&rating.CreatedAt,
			&rating.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func GetRatingTransaction(pool *pgxpool.Pool, transactionID int) ([]models.Rating, error) {

	rows, err := pool.Query(context.Background(), `
		SELECT r.id, r.users_id, r.products_id, r.transactions_id, r.rating, r.review, r.created_at, r.updated_at
		FROM ratings r
		WHERE r.transactions_id = $1
		ORDER BY r.created_at DESC
	`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []models.Rating
	for rows.Next() {
		var rating models.Rating
		err := rows.Scan(
			&rating.ID,
			&rating.UsersID,
			&rating.ProductsID,
			&rating.TransactionsID,
			&rating.Rating,
			&rating.Review,
			&rating.CreatedAt,
			&rating.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func UpdateRating(pool *pgxpool.Pool, rating *models.Rating) error {
	_, err := pool.Exec(
		context.Background(),
		`
		UPDATE ratings
		SET rating = $1, review = $2, updated_at = $3
		WHERE id = $4 AND users_id = $5
	`,
		rating.Rating,
		rating.Review,
		time.Now(),
		rating.ID,
		rating.UsersID,
	)

	if err != nil {
		return err
	}

	return nil
}

func DeleteRatinh(pool *pgxpool.Pool, id int) error {
	_, err := pool.Exec(context.Background(), `DELETE FROM ratings WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func CheckUserPurchased(pool *pgxpool.Pool, userID, productID, transactionID int) (bool, error) {

	var exists bool
	err := pool.QueryRow(context.Background(), `
		SELECT EXISTS(
			SELECT 1 FROM transaction_items ti
			INNER JOIN transactions t ON ti.transactions_id = t.id
			WHERE t.users_id = $1 
			AND ti.products_id = $2 
			AND t.id = $3
		)
	`, userID, productID, transactionID).Scan(&exists)
	return exists, err
}

func CheckAlreadyRated(pool *pgxpool.Pool, userID, productID, transactionID int) (bool, error) {

	var exists bool
	err := pool.QueryRow(context.Background(), `
		SELECT EXISTS(
			SELECT 1 FROM ratings
			WHERE users_id = $1 
			AND products_id = $2 
			AND transactions_id = $3
		)
	`, userID, productID, transactionID).Scan(&exists)
	return exists, err
}

func GetUnratedProductsByTransaction(pool *pgxpool.Pool, userID, transactionID int) ([]models.UnratedProduct, error) {
	rows, err := pool.Query(context.Background(), `
		SELECT DISTINCT ti.products_id, p.name, ti.quantity
		FROM transaction_items ti
		INNER JOIN transactions t ON ti.transactions_id = t.id
		INNER JOIN products p ON ti.products_id = p.id
		LEFT JOIN ratings r ON r.products_id = ti.products_id 
			AND r.transactions_id = t.id 
			AND r.users_id = t.users_id
		WHERE t.id = $1 
		AND t.users_id = $2
		AND r.id IS NULL
		ORDER BY p.name
	`, transactionID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.UnratedProduct
	for rows.Next() {
		var product models.UnratedProduct
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
