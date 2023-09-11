package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
)

const (
	aboutUsSignature = "App\\Models\\AboutUs"
	mediaTable       = "media"
)

type MediaRepository struct {
	client *sql.DB
	logger log.Logger
}

func NewMediaRepository(client *sql.DB, logger log.Logger) *MediaRepository {
	return &MediaRepository{client: client, logger: logger}
}

func (m MediaRepository) ListAboutUs(ctx context.Context) ([]entities.Media, error) {
	query, args, err := m.makeListAboutUSQuery()
	m.logger.Log("query", query, "args", fmt.Sprintf("%v", args))
	if err != nil {
		return nil, err
	}

	rows, err := m.client.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	mediaList := []entities.Media{}
	for rows.Next() {
		m.logger.Log("row", "ssss")
		var media entities.Media
		scans := []interface{}{
			&media.ModelType,
			&media.ModelID,
			&media.FileName,
			&media.OrderColumn,
			&media.MediaType,
		}
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		mediaList = append(mediaList, media)
	}
	m.logger.Log("len", fmt.Sprintf("%v", mediaList))
	return mediaList, nil
}

func (m MediaRepository) makeListAboutUSQuery() (string, []interface{}, error) {
	return sq.Select(
		"media.model_type",
		"media.model_id",
		"media.file_name",
		"media.order_column",
		"about_us.type",
	).
		From(mediaTable).
		Join("about_us ON about_us.id=media.model_id").
		Where(sq.Eq{"media.model_type": aboutUsSignature}).
		Where(sq.Eq{"media.custom_properties->\"$.default\"": 1}).
		ToSql()
}
