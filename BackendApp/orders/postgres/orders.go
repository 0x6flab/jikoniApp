package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/0x6flab/jikoniApp/BackendApp/orders"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
)

var _ orders.OrderRepository = (*orderRepo)(nil)

type orderRepo struct {
	db *sqlx.DB
}

// NewOrderRepo instantiates a PostgreSQL
// implementation of Users repository.
func NewOrderRepo(db *sqlx.DB) orders.OrderRepository {
	return &orderRepo{
		db: db,
	}
}

func (repo orderRepo) Save(ctx context.Context, order orders.Order) (string, error) {
	q := `INSERT INTO orders (id, vendor, name, price, place, status, metadata, created_at, updated_at)
		  VALUES (:id, :vendor, :name, :price, :place, :status, :metadata, :created_at, :updated_at) RETURNING id`

	dbo, err := toDBOrder(order)
	if err != nil {
		return "", multierr.Combine(errors.ErrCreateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbo)
	if err != nil {
		return "", handleError(err, errors.ErrCreateEntity)
	}
	defer row.Close()
	row.Next()
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (repo orderRepo) RetrieveByID(ctx context.Context, id string) (orders.Order, error) {
	q := `SELECT id, vendor, name, price, place, status, metadata, created_at, updated_at FROM orders WHERE id = $1`

	dbc := dbOrder{
		ID: id,
	}
	if err := repo.db.QueryRowxContext(ctx, q, id).StructScan(&dbc); err != nil {
		if err == sql.ErrNoRows {
			return orders.Order{}, multierr.Combine(errors.ErrNotFound, err)
		}
		return orders.Order{}, multierr.Combine(errors.ErrViewEntity, err)
	}
	return toOrder(dbc)
}

func (repo orderRepo) RetrieveAll(ctx context.Context, pm orders.PageMetadata) (orders.OrdersPage, error) {
	mq, mp, err := createMetadataQuery("", pm.Metadata)
	if err != nil {
		return orders.OrdersPage{}, multierr.Combine(errors.ErrViewEntity, err)
	}
	var query []string
	var emq string
	if mq != "" {
		query = append(query, mq)
	}
	if pm.Vendor != "" {
		query = append(query, fmt.Sprintf("vendor = '%s'", pm.Vendor))
	}
	if pm.Name != "" {
		query = append(query, fmt.Sprintf("name = '%s'", pm.Name))
	}
	if pm.Price != 0 {
		query = append(query, fmt.Sprintf("price = '%d'", pm.Price))
	}
	if pm.Place != "" {
		query = append(query, fmt.Sprintf("place = '%s'", pm.Place))
	}
	if pm.Status != "" {
		query = append(query, fmt.Sprintf("status = '%s'", pm.Status))
	}
	if len(query) > 0 {
		emq = fmt.Sprintf(" WHERE %s", strings.Join(query, " AND "))
	}

	q := fmt.Sprintf(`SELECT id, vendor, name, price, place, status, metadata, created_at, updated_at FROM orders %s ORDER BY created_at LIMIT :limit OFFSET :offset;`, emq)
	params := map[string]interface{}{
		"limit":    pm.Limit,
		"offset":   pm.Offset,
		"metadata": mp,
	}
	rows, err := repo.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return orders.OrdersPage{}, multierr.Combine(errors.ErrViewEntity, err)
	}
	defer rows.Close()
	var items []orders.Order
	for rows.Next() {
		dbc := dbOrder{}
		if err := rows.StructScan(&dbc); err != nil {
			return orders.OrdersPage{}, multierr.Combine(errors.ErrViewEntity, err)
		}
		c, err := toOrder(dbc)
		if err != nil {
			return orders.OrdersPage{}, err
		}
		items = append(items, c)
	}
	cq := fmt.Sprintf(`SELECT COUNT(*) FROM orders %s;`, emq)

	total, err := total(ctx, repo.db, cq, params)
	if err != nil {
		return orders.OrdersPage{}, multierr.Combine(errors.ErrViewEntity, err)
	}
	page := orders.OrdersPage{
		Orders: items,
		PageMetadata: orders.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}
	return page, nil
}

func (repo orderRepo) Update(ctx context.Context, order orders.Order) (string, error) {
	var query []string
	var upq string
	if order.Vendor != "" {
		query = append(query, "vendor = :vendor,")
	}
	if order.Name != "" {
		query = append(query, "name = :name,")
	}
	if order.Price != 0 {
		query = append(query, "price = :price,")
	}
	if order.Place != "" {
		query = append(query, "place = :place,")
	}
	if order.Status != "" {
		query = append(query, "status = :status,")
	}
	if order.Metadata != nil {
		query = append(query, "metadata = :metadata,")
	}
	if len(query) > 0 {
		upq = strings.Join(query, " ")
	}
	q := fmt.Sprintf(`UPDATE orders SET %s updated_at = :updated_at WHERE id = :id RETURNING id`, upq)

	dbu, err := toDBOrder(order)
	if err != nil {
		return "", multierr.Combine(errors.ErrUpdateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbu)
	if err != nil {
		return "", multierr.Combine(errors.ErrUpdateEntity, err)
	}
	defer row.Close()
	row.Next()
	var id string
	if err := row.Scan(&id); err != nil {
		return "", multierr.Combine(errors.ErrUpdateEntity, err)
	}
	return id, nil
}

func (repo orderRepo) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM orders WHERE id = :id`

	dbu := dbOrder{
		ID: id,
	}
	if _, err := repo.db.NamedExecContext(ctx, q, dbu); err != nil {
		return multierr.Combine(errors.ErrRemoveEntity, err)
	}
	return nil
}

func total(ctx context.Context, db *sqlx.DB, query string, params interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, nil
}

type dbOrder struct {
	ID        string    `db:"id,omitempty"`
	Vendor    string    `db:"vendor,omitempty"`
	Name      string    `db:"name,omitempty"`
	Price     uint64    `db:"price,omitempty"`
	Place     string    `db:"place,omitempty"`
	Metadata  []byte    `db:"metadata,omitempty"`
	Status    string    `db:"status,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
}

func toDBOrder(order orders.Order) (dbOrder, error) {
	data := []byte("{}")
	if len(order.Metadata) > 0 {
		b, err := json.Marshal(order.Metadata)
		if err != nil {
			return dbOrder{}, multierr.Combine(errors.ErrMalformedEntity, err)
		}
		data = b
	}
	return dbOrder{
		ID:        order.ID,
		Vendor:    order.Vendor,
		Name:      order.Name,
		Price:     order.Price,
		Place:     order.Place,
		Metadata:  data,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}, nil
}

func toOrder(order dbOrder) (orders.Order, error) {
	var metadata map[string]interface{}
	if order.Metadata != nil {
		if err := json.Unmarshal([]byte(order.Metadata), &metadata); err != nil {
			return orders.Order{}, multierr.Combine(errors.ErrMalformedEntity, err)
		}
	}
	return orders.Order{
		ID:        order.ID,
		Vendor:    order.Vendor,
		Name:      order.Name,
		Price:     order.Price,
		Place:     order.Place,
		Metadata:  metadata,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}, nil
}

func createMetadataQuery(entity string, um orders.Metadata) (string, []byte, error) {
	if len(um) == 0 {
		return "", nil, nil
	}
	param, err := json.Marshal(um)
	if err != nil {
		return "", nil, err
	}
	query := fmt.Sprintf("%smetadata @> :metadata", entity)
	return query, param, nil
}
