package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sixojke/test-service/internal/domain"
)

type GoodsPostgres struct {
	db *sqlx.DB
}

func NewGoodsPostgres(db *sqlx.DB) *GoodsPostgres {
	return &GoodsPostgres{
		db: db,
	}
}

func (r *GoodsPostgres) Create(inp domain.ItemCreateInp) (*domain.Item, error) {
	var desc sql.NullString

	var item domain.Item
	query := fmt.Sprintf(`INSERT INTO %s (project_id, name) VALUES ($1, $2)
		RETURNING id, project_id, name, description, priority, removed, created_at`, goods)

	if err := r.db.QueryRow(query, inp.ProjectId, inp.Name).Scan(&item.Id, &item.ProjectId,
		&item.Name, &desc, &item.Priority, &item.Removed, &item.CreatedAt); err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}

	if desc.Valid {
		item.Description = desc.String
	} else {
		item.Description = ""
	}

	return &item, nil
}

func (r *GoodsPostgres) Update(inp domain.ItemUpdateInp) (*domain.Item, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, inp.Name)
		argId++
	}

	if inp.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, inp.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	var item domain.Item
	query := fmt.Sprintf(`UPDATE %s SET %s where id='%v' AND project_id='%v';`,
		goods, setQuery, inp.Id, inp.ProjectId)
	if err := txUpdate(r.db, query, args, inp.Id); err != nil {
		return nil, fmt.Errorf("tx: %v", err)
	}

	query = fmt.Sprintf(`SELECT id, project_id, name, COALESCE(description, '') AS description, priority, removed, created_at 
		FROM %v WHERE id = $1`, goods)
	if err := r.db.QueryRow(query, inp.Id).Scan(&item.Id, &item.ProjectId, &item.Name, &item.Description,
		&item.Priority, &item.Removed, &item.CreatedAt); err != nil {
		return nil, fmt.Errorf("select query: %v", err)
	}

	return &item, nil
}

func txUpdate(db *sqlx.DB, query string, args []interface{}, id int) error {
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("create transaction: %v", err)
	}

	if _, err := tx.Exec(fmt.Sprintf("SELECT * FROM %s WHERE id = $1 FOR UPDATE", goods), id); err != nil {
		tx.Rollback()
		return fmt.Errorf("select for update: %v", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		tx.Rollback()
		return fmt.Errorf("query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %v", err)
	}

	return nil
}

func (r *GoodsPostgres) Delete(inp domain.ItemDeleteInp) (*domain.ItemDeleteOut, error) {
	var item domain.ItemDeleteOut
	args := []interface{}{inp.Id, inp.ProjectId}

	var isRemoved bool
	query := fmt.Sprintf("SELECT removed FROM %s WHERE id = $1", goods)
	if err := r.db.QueryRow(query, inp.Id).Scan(&isRemoved); err != nil {
		return nil, fmt.Errorf("error checking if item is removed")
	}
	if isRemoved {
		return nil, fmt.Errorf("item is already removed")
	}

	query = fmt.Sprintf("UPDATE %s SET removed = true WHERE id = $1 AND project_id = $2", goods)
	if err := txUpdate(r.db, query, args, inp.Id); err != nil {
		return nil, fmt.Errorf("tx: %v", err)
	}

	query = fmt.Sprintf("SELECT id, project_id, removed FROM %s WHERE id = $1", goods)
	if err := r.db.QueryRow(query, inp.Id).Scan(&item.Id, &item.ProjectId, &item.Removed); err != nil {
		return nil, fmt.Errorf("record not found")
	}

	return &item, nil
}

func (r *GoodsPostgres) GetList(limit, offset int) (*domain.List, error) {
	var g []*domain.Item
	var total, removedCount int

	totalQuery := fmt.Sprintf("SELECT COUNT(*) FROM %v;", goods)
	err := r.db.Get(&total, totalQuery)
	if err != nil {
		return nil, fmt.Errorf("query total count: %v", err)
	}

	removedQuery := fmt.Sprintf("SELECT COUNT(*) FROM %v WHERE removed = true;", goods)
	err = r.db.Get(&removedCount, removedQuery)
	if err != nil {
		return nil, fmt.Errorf("query removed count: %v", err)
	}

	query := fmt.Sprintf(`SELECT id, project_id, name, priority, removed, created_at, 
        COALESCE(description, '') AS description FROM %v WHERE removed = FALSE LIMIT $1 OFFSET $2;`, goods)

	if err := r.db.Select(&g, query, limit, offset); err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}

	return &domain.List{
		Total:   total,
		Removed: removedCount,
		Limit:   limit,
		Offset:  offset,
		Goods:   g,
	}, nil
}

func (r *GoodsPostgres) GetById(id int) (*domain.Item, error) {
	var item domain.Item
	query := fmt.Sprintf(`SELECT id, project_id, name, COALESCE(description, '') AS description,
		priority, removed FROM %v where id = $1`, goods)

	if err := r.db.Get(&item, query, id); err != nil {
		return nil, fmt.Errorf("get by id query: %v", err)
	}

	return &item, nil
}

func (r *GoodsPostgres) Reprioritiize(inp domain.ItemReprioritiizeInp) ([]domain.ItemReprioritiizeOut, error) {
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create transaction: %v", err)
	}

	if _, err := tx.Exec(fmt.Sprintf("SELECT * FROM %s WHERE id = $1 FOR UPDATE", goods), inp.Id); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("select for update: %v", err)
	}

	query := fmt.Sprintf(`UPDATE %v SET priority = $2 WHERE id = $1;`, goods)
	if _, err := tx.Exec(query, inp.Id, inp.Priority); err != nil {
		return nil, fmt.Errorf("error updating priority for item: %v", err)
	}

	query = fmt.Sprintf(`UPDATE %v SET priority = priority + 1 WHERE priority >= $1 AND NOT id = $2;`, goods)
	_, err = tx.Exec(query, inp.Priority, inp.Id)
	if err != nil {
		return nil, fmt.Errorf("error updating priorities for items with higher priority: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %v", err)
	}

	var g []domain.ItemReprioritiizeOut
	query = fmt.Sprintf(`SELECT id, priority FROM %v WHERE priority >= $1 ORDER BY priority;`, goods)
	if err := r.db.Select(&g, query, inp.Priority); err != nil {
		return nil, fmt.Errorf("error fetching updated items: %v", err)
	}

	return g, nil
}
