package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/0x6flab/jikoniApp/BackendApp/shops"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
)

var _ shops.ShopRepository = (*shopRepo)(nil)

type shopRepo struct {
	db *sqlx.DB
}

// NewShopRepo instantiates a PostgreSQL
// implementation of Users repository.
func NewShopRepo(db *sqlx.DB) shops.ShopRepository {
	return &shopRepo{
		db: db,
	}
}

func (repo shopRepo) Save(ctx context.Context, shop shops.Shop) (string, error) {
	q := `INSERT INTO shops (id, name, email, number, metadata, created_at, updated_at)
		  VALUES (:id, :name, :email, :number, :metadata, :created_at, :updated_at) RETURNING id`

	dbs, err := toDBShop(shop)
	if err != nil {
		return "", multierr.Combine(errors.ErrCreateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbs)
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

func (repo shopRepo) RetrieveByID(ctx context.Context, id string) (shops.Shop, error) {
	q := `SELECT id, name, email, number, metadata, created_at, updated_at FROM shops WHERE id = $1`

	dbs := dbShop{
		ID: id,
	}
	if err := repo.db.QueryRowxContext(ctx, q, id).StructScan(&dbs); err != nil {
		if err == sql.ErrNoRows {
			return shops.Shop{}, multierr.Combine(errors.ErrNotFound, err)
		}
		return shops.Shop{}, multierr.Combine(errors.ErrViewEntity, err)
	}
	return toShop(dbs)
}

func (repo shopRepo) RetrieveAll(ctx context.Context, pm shops.PageMetadata) (shops.ShopsPage, error) {
	mq, mp, err := createMetadataQuery("", pm.Metadata)
	if err != nil {
		return shops.ShopsPage{}, multierr.Combine(errors.ErrViewEntity, err)
	}
	var query []string
	var emq string
	if mq != "" {
		query = append(query, mq)
	}
	if pm.Name != "" {
		query = append(query, fmt.Sprintf("name = '%s'", pm.Name))
	}
	if pm.Email != "" {
		query = append(query, fmt.Sprintf("email = '%s'", pm.Email))
	}
	if pm.Number != "" {
		query = append(query, fmt.Sprintf("number = '%s'", pm.Number))
	}
	if len(query) > 0 {
		emq = fmt.Sprintf(" WHERE %s", strings.Join(query, " AND "))
	}

	q := fmt.Sprintf(`SELECT id, name, email, number, metadata, created_at, updated_at FROM shops %s ORDER BY created_at LIMIT :limit OFFSET :offset;`, emq)
	params := map[string]interface{}{
		"limit":    pm.Limit,
		"offset":   pm.Offset,
		"metadata": mp,
	}
	rows, err := repo.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return shops.ShopsPage{}, multierr.Combine(errors.ErrViewEntity, err)
	}
	defer rows.Close()
	var items []shops.Shop
	for rows.Next() {
		dbs := dbShop{}
		if err := rows.StructScan(&dbs); err != nil {
			return shops.ShopsPage{}, multierr.Combine(errors.ErrViewEntity, err)
		}
		s, err := toShop(dbs)
		if err != nil {
			return shops.ShopsPage{}, err
		}
		items = append(items, s)
	}
	cq := fmt.Sprintf(`SELECT COUNT(*) FROM shops %s;`, emq)

	total, err := total(ctx, repo.db, cq, params)
	if err != nil {
		return shops.ShopsPage{}, multierr.Combine(errors.ErrViewEntity, err)
	}
	page := shops.ShopsPage{
		Shops: items,
		PageMetadata: shops.PageMetadata{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}
	return page, nil
}

func (repo shopRepo) Update(ctx context.Context, shop shops.Shop) (string, error) {
	var query []string
	var upq string
	if shop.Name != "" {
		query = append(query, "name = :name,")
	}
	if shop.Email != "" {
		query = append(query, "email = :email,")
	}
	if shop.Number != "" {
		query = append(query, "number = :number,")
	}
	if shop.Metadata != nil {
		query = append(query, "metadata = :metadata,")
	}
	if len(query) > 0 {
		upq = strings.Join(query, " ")
	}
	q := fmt.Sprintf(`UPDATE shops SET %s updated_at = :updated_at WHERE id = :id RETURNING id`, upq)

	dbs, err := toDBShop(shop)
	if err != nil {
		return "", multierr.Combine(errors.ErrUpdateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbs)
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

func (repo shopRepo) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM shops WHERE id = :id`

	dbs := dbShop{
		ID: id,
	}
	if _, err := repo.db.NamedExecContext(ctx, q, dbs); err != nil {
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

type dbShop struct {
	ID        string    `db:"id,omitempty"`
	Name      string    `db:"name,omitempty"`
	Email     string    `db:"email,omitempty"`
	Number    string    `db:"number,omitempty"`
	Metadata  []byte    `db:"metadata,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
}

func toDBShop(shop shops.Shop) (dbShop, error) {
	data := []byte("{}")
	if len(shop.Metadata) > 0 {
		b, err := json.Marshal(shop.Metadata)
		if err != nil {
			return dbShop{}, multierr.Combine(errors.ErrMalformedEntity, err)
		}
		data = b
	}
	return dbShop{
		ID:        shop.ID,
		Name:      shop.Name,
		Email:     shop.Email,
		Number:    shop.Number,
		Metadata:  data,
		CreatedAt: shop.CreatedAt,
		UpdatedAt: shop.UpdatedAt,
	}, nil
}

func toShop(shop dbShop) (shops.Shop, error) {
	var metadata map[string]interface{}
	if shop.Metadata != nil {
		if err := json.Unmarshal([]byte(shop.Metadata), &metadata); err != nil {
			return shops.Shop{}, multierr.Combine(errors.ErrMalformedEntity, err)
		}
	}
	return shops.Shop{
		ID:        shop.ID,
		Name:      shop.Name,
		Email:     shop.Email,
		Number:    shop.Number,
		Metadata:  metadata,
		CreatedAt: shop.CreatedAt,
		UpdatedAt: shop.UpdatedAt,
	}, nil
}

func createMetadataQuery(entity string, um shops.Metadata) (string, []byte, error) {
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
