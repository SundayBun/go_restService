package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"goWebService/models"
)

type Repository interface {
	Save(ctx context.Context, account *models.AccountModel) (*models.AccountModel, error)
	Update(ctx context.Context, account *models.AccountModel) (*models.AccountModel, error)
	DeleteById(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*models.AccountModel, error)
}

type pgRepository struct {
	db *sqlx.DB
}

func NewPgRepository(db *sqlx.DB) Repository {
	return &pgRepository{db: db}
}

func (p *pgRepository) Save(ctx context.Context, account *models.AccountModel) (*models.AccountModel, error) {
	var a models.AccountModel
	if err := p.db.QueryRowxContext(
		ctx,
		save,
		&account.FirstName,
		&account.LastName,
	).StructScan(&a); err != nil {
		return nil, errors.Wrap(err, "pgRepository.Save.QueryRowxContext")
	}

	return &a, nil
}

func (p *pgRepository) Update(ctx context.Context, account *models.AccountModel) (*models.AccountModel, error) {
	var a models.AccountModel
	if err := p.db.QueryRowxContext(
		ctx,
		update,
		&account.FirstName,
		&account.LastName,
		&account.Id,
	).StructScan(&a); err != nil {
		return nil, errors.Wrap(err, "pgRepository.Update.QueryRowxContext")
	}

	return &a, nil
}

func (p *pgRepository) DeleteById(ctx context.Context, id int) error {
	result, err := p.db.ExecContext(
		ctx,
		deleteById,
		id,
	)

	if err != nil {
		return errors.Wrap(err, "pgRepository.Delete.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "pgRepository.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "pgRepository.Delete.rowsAffected")
	}

	return nil
}

func (p *pgRepository) GetById(ctx context.Context, id int) (*models.AccountModel, error) {
	var a = &models.AccountModel{}
	if err := p.db.GetContext(
		ctx, a,
		findById,
		id,
	); err != nil {
		return nil, errors.Wrap(err, "pgRepository.GetById.GetContext")
	}
	return a, nil
}
