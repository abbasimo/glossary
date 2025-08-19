package terms

import (
	"context"
	"database/sql"
)

type Term struct {
	ID           int
	TermEn       string
	DefinitionEn string
	Context      string
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) All(ctx context.Context) ([]Term, error) {
	rows, err := r.db.QueryContext(ctx, "select id, term_en, definition_en, context from terms order by id desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var terms []Term
	for rows.Next() {
		var t Term
		if err := rows.Scan(&t.ID, &t.TermEn, &t.DefinitionEn, &t.Context); err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	return terms, nil
}

func (r *Repository) ByID(ctx context.Context, id int) (*Term, error) {
	row := r.db.QueryRowContext(ctx, "select id, term_en, definition_en, context from terms where id=$1", id)
	var t Term
	if err := row.Scan(&t.ID, &t.TermEn, &t.DefinitionEn, &t.Context); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) Create(ctx context.Context, t *Term) error {
	return r.db.QueryRowContext(ctx,
		"insert into terms (term_en, definition_en, context) values ($1,$2,$3) returning id",
		t.TermEn, t.DefinitionEn, t.Context,
	).Scan(&t.ID)
}
func (r *Repository) Update(ctx context.Context, t *Term) error {
	_, err := r.db.ExecContext(ctx,
		"update terms set term_en=$1, definition_en=$2, context=$3 where id=$4",
		t.TermEn, t.DefinitionEn, t.Context, t.ID)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "delete from terms where id=$1", id)
	return err
}
