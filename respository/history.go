package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetHistorys(pool *pgxpool.Pool, userID int, page int, limit int, month string, status string) (models.PaginationResponseHistory, error){
	var historys []models.History
	offset := (page - 1)* limit

	query := `SELECT t.id, t.invoice_num, t.created_at, ts.status, t.total, MAX(pi.image) AS image
	FROM transactions t
	LEFT JOIN status_transactions ts ON ts.id = t.status_transactions_id
	LEFT JOIN transaction_items ti ON ti.transactions_id = t.id
	LEFT JOIN products p ON ti.products_id = p.id
	LEFT JOIN product_images pi ON pi.products_id = p.id
	WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if userID != 0 {
		query += fmt.Sprintf(" AND t.users_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	if month != ""{
		query += fmt.Sprintf(" AND EXTRACT(MONTH FROM t.created_at) = $%d", argIndex)
		args = append(args, month)
		argIndex++
	}

	if status != ""{
		query+= fmt.Sprintf(" AND t.status_transactions_id = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	query += `
	GROUP BY t.id, t.invoice_num, t.created_at, ts.status, t.total
	ORDER BY t.created_at DESC
	OFFSET $` + fmt.Sprint(argIndex) + ` LIMIT $` + fmt.Sprint(argIndex+1)
	args = append(args, offset, limit)

	rows, err := pool.Query(context.Background(), query, args...)
	if err != nil {
		fmt.Println("Error: Failed to get data transaction", err)
	}

	defer rows.Close()

	for rows.Next(){
		var item models.History

		err := rows.Scan(
			&item.Id,
			&item.Invoice,
			&item.Date,
			&item.Status,
			&item.Total,
			&item.ImageProduct,
		)

		if err != nil {
			fmt.Println("Error: Failed to scanning", err)
		}

		historys = append(historys, item)
	}

	countQuery := `SELECT COUNT(*) FROM transactions t WHERE 1=1`

	countArgs := []interface{}{}
	countIndex := 1

	if userID != 0 {
		countQuery += fmt.Sprintf(" AND t.users_id = $%d", countIndex)
		countArgs = append(countArgs, userID)
		countIndex++
	}

	if month != "" {
		countQuery += fmt.Sprintf(" AND EXTRACT(MONTH FROM t.created_at) = $%d", countIndex)
		countArgs = append(countArgs, month)
		countIndex++
	}
	if status != "" {
		countQuery += fmt.Sprintf(" AND t.status_transactions_id = $%d", countIndex)
		countArgs = append(countArgs, status)
		countIndex++
	}

	var total int
	err = pool.QueryRow(context.Background(), countQuery, countArgs...).Scan(&total)
	if err != nil {
		fmt.Println("Error counting transactions:", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	links := make(map[string]string)
	if page > 1 {
		links["prev"] = fmt.Sprintf("/history?page=%d", page-1)
	} else {
		links["prev"] = "null"
	}
	if page < totalPages {
		links["next"] = fmt.Sprintf("/history?page=%d", page+1)
	} else {
		links["next"] = "null"
	}

	response := models.PaginationResponseHistory{
		Data: historys,
		Page: page,
		Limit: limit,
		Total: totalPages,
		Links: links,
	}

	return response, nil
}