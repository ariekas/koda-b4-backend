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

func DetailHistory(pool *pgxpool.Pool, userID int, historyId int) (models.DetailHistory, error) {
	var historys models.DetailHistory

	err := pool.QueryRow(context.Background(), `
	SELECT 
	t.id,
    t.invoice_num,
    t.name_user,
    t.address_user,
    t.phone_user,
    pm.name AS payment_method,
    d.name AS delivery_method,
    st.status,
    t.total,
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'item_id', ti.id,
           'image', pi.image,
		   'name', p.name,
		   'price', p.price,
		   'price discount', p.price_discounts,
            'price discount', p.price_discounts,
			'quantity', ti.quantity,
            'size', sp.name,
            'variant', vp.name,
            'subtotal', ti.subtotal
        )
    ) AS items
FROM transactions t
LEFT JOIN payment_methods pm ON t.payment_methods_id = pm.id
LEFT JOIN status_transactions st ON t.status_transactions_id = st.id
LEFT JOIN transaction_items ti ON ti.transactions_id = t.id
LEFT JOIN products p ON ti.products_id = p.id
LEFT JOIN ( SELECT DISTINCT ON (products_id) products_id, image
FROM product_images
ORDER BY products_id, id ASC
) pi ON pi.products_id = p.id
LEFT JOIN size_products sp ON ti.size_products_id = sp.id
LEFT JOIN variant_products vp ON ti.variant_products_id = vp.id
LEFT JOIN deliverys d ON t.deliverys_id = d.id
WHERE t.id = $1 AND t.users_id = $2
GROUP BY 
    t.id, t.invoice_num, t.name_user, t.address_user, t.phone_user, 
    pm.name, d.name, st.status, t.total, p.price_discounts`,historyId, userID).Scan(
		&historys.Id, &historys.Invoice, &historys.Fullname, &historys.Address, &historys.Phone, &historys.Payment, &historys.Delivery, &historys.Status, &historys.Total, &historys.HistoryItems)

		return historys, err
}