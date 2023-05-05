package postgresql

import (
	"app/api/models"
	"fmt"
	"time"

	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type examPepo struct {
	db *pgxpool.Pool
}

func NewExamRepo(db *pgxpool.Pool) *examPepo {
	return &examPepo{
		db: db,
	}
}

func (r *examPepo) Create(ctx context.Context, req *models.CreatePromoCode) (int, error) {

	query := `INSERT INTO promocode("id","name",
		"discount",
		"discount_type",
		"order_limit_price"
		) 
		VALUES((SELECT MAX(id) + 1 FROM promocode),$1,$2,$3,$4) RETURNING id`
	id := 0
	err := r.db.QueryRow(ctx, query, req.Name, req.Discount, req.DiscountType, req.OrderLimitPrice).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *examPepo) GetByID(ctx context.Context, req *models.PromocodePrimaryKey) (*models.Promocode, error) {
	var (
		query     string
		promocode models.Promocode
	)

	query = `SELECT * FROM promocode WHERE id =$1`

	err := r.db.QueryRow(ctx, query, req.PromocodeId).Scan(
		&promocode.Id, &promocode.Name, &promocode.Discount, &promocode.DiscountType, &promocode.OrderLimitPrice)
	if err != nil {
		return nil, err
	}

	return &promocode, nil
}

func (r *examPepo) GetList(ctx context.Context, req *models.GetListBrandRequest) (resp *models.GetListPromocodeResponse, err error) {

	resp = &models.GetListPromocodeResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id, 
			name, 
			discount,
			discount_type,
			order_limit_price
		FROM promocode
	`

	if len(req.Search) > 0 {
		filter += " AND name ILIKE '%' || '" + req.Search + "' || '%' "
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var promocode models.Promocode
		err = rows.Scan(
			&resp.Count,
			&promocode.Id,
			&promocode.Name,
			&promocode.Discount,
			&promocode.DiscountType,
			&promocode.OrderLimitPrice,
		)
		if err != nil {
			return nil, err
		}

		resp.Promocodes = append(resp.Promocodes, &promocode)
	}

	return resp, nil
}

func (r *examPepo) Delete(ctx context.Context, req *models.PromocodePrimaryKey) (int64, error) {
	query := `
		DELETE 
		FROM promocode
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, req.PromocodeId)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *examPepo) TotalSumm(ctx context.Context, req *models.OrderItemPrimaryKey) (float64, error) {
	query := `
		SELECT SUM(products.list_price), 
		       COALESCE(CASE WHEN promocode.discount_type = 'percent' THEN promocode.discount * SUM(products.list_price) / 100.0 ELSE promocode.discount END, 0) as discount
		FROM orders
		LEFT JOIN order_items ON orders.order_id = order_items.order_id
		LEFT JOIN products ON order_items.product_id = products.product_id
		LEFT JOIN promocode ON orders.promocode_id = promocode.id 
		WHERE orders.order_id = $1
		GROUP BY orders.order_id, promocode.discount, promocode.discount_type`

	var summ float64
	var discount float64

	err := r.db.QueryRow(ctx, query, req.OrderId).Scan(&summ, &discount)
	if err != nil {
		return 0, err
	}
	// fmt.Println(r.Report("22-01-2004"))
	return summ - discount, nil
}

func (r *examPepo) Report(ctx context.Context, req *models.Date) (res []models.StaffDate, err error) {
	query := `SELECT
    staffs.first_name || ' ' || staffs.last_name AS "employe",  categories.category_name AS "category",
       products.product_name AS "product",   order_items.quantity AS "quantity",   order_items.list_price * order_items.quantity AS "summ"
FROM orders
         JOIN order_items ON orders.order_id = order_items.order_id
         JOIN products ON order_items.product_id = products.product_id
         JOIN categories ON products.category_id = categories.category_id
         JOIN staffs ON orders.staff_id = staffs.staff_id
WHERE orders.order_date = $1`

	var year string

	if req.Day == "" {
		dt := time.Now()
		year = dt.Format("2006-01-02")
	} else {
		year = req.Day
	}

	date, error := time.Parse("2006-01-02", year)
	if error != nil {
		fmt.Println(error)
		return
	}

	rows, err := r.db.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var s models.StaffDate
		err = rows.Scan(
			&s.StaffName,
			&s.Category,
			&s.Product,
			&s.Quantity,
			&s.Summ,
		)
		res = append(res, s)
		if err != nil {
			return res, err
		}
	}
	return res, nil

}