package repo

import (
	"context"

	"T4_test_case/internal/restserver/model"
	"T4_test_case/pkg/db"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type FileRepository interface {
	SaveFile(ctx context.Context, file model.File) (int64, error)
	GetFile(ctx context.Context, id int64) (*model.File, error)
}

type fileRepository struct {
	store db.Store
}

// SaveFile - save file to db
func (f *fileRepository) SaveFile(ctx context.Context, file model.File) (int64, error) {
	sql, args, err := sq.Insert(file.Table()).
		PlaceholderFormat(sq.Dollar).
		Columns(file.InsertColumns()...).
		Values(file.InsertValues()...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "[file-repo] unable to build INSERT query")
	}

	id, err := f.store.Insert(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "[file-repo] f.store.Insert()")
	}

	return id, nil
}

// GetFile - get file by id from db
func (f *fileRepository) GetFile(ctx context.Context, id int64) (*model.File, error) {
	file := model.File{}
	sql, args, err := sq.Select("*").
		From(file.Table()).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[file-repo] unable to build SELECT query")
	}

	err = f.store.Get(ctx, &file, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "[file-repo] f.store.Get()")
	}

	return &file, nil
}

func NewFileRepository(store db.Store) FileRepository {
	return &fileRepository{
		store: store,
	}
}
