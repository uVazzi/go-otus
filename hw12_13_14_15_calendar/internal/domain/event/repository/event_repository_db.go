package eventrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/txrunner"
)

type repositoryDB struct {
	db *sql.DB
}

func ProvideRepositoryDB(db *sql.DB) Repository {
	return &repositoryDB{
		db: db,
	}
}

func (repo *repositoryDB) Save(ctx context.Context, tx txrunner.Tx, eventModel *event.Event) error {
	if eventModel == nil {
		return ErrEmptyEvent
	}

	if eventModel.UUID == uuid.Nil {
		eventModel.UUID = uuid.New()
	}

	const query = `
		INSERT INTO event 
		    (uuid, user_uuid, title, description, start_date, end_date, delay_notification, delay_notification_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (uuid)
		DO UPDATE SET
			user_uuid = EXCLUDED.user_uuid,
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date,
			delay_notification = EXCLUDED.delay_notification,
			delay_notification_type = EXCLUDED.delay_notification_type
		RETURNING uuid
	`

	return tx.QueryRowContext(ctx, query,
		eventModel.UUID,
		eventModel.UserUUID,
		eventModel.Title,
		eventModel.Description,
		eventModel.StartDate,
		eventModel.EndDate,
		eventModel.DelayNotification,
		eventModel.DelayNotificationType,
	).Scan(&eventModel.UUID)
}

func (repo *repositoryDB) Delete(ctx context.Context, tx txrunner.Tx, eventUUID uuid.UUID) error {
	result, err := tx.ExecContext(ctx, `DELETE FROM event WHERE uuid = $1`, eventUUID)
	if err != nil {
		return err
	}
	totalDeleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if totalDeleted == 0 {
		return ErrEventNotFoundByDelete
	}

	return nil
}

func (repo *repositoryDB) Get(ctx context.Context, filter Filter) (*event.Event, error) {
	err := CheckFilter(filter)
	if err != nil {
		return nil, err
	}

	query, args := buildQuery(filter, true)
	row := repo.db.QueryRowContext(ctx, query, args...)

	var eventModel event.Event
	err = row.Scan(
		&eventModel.UUID,
		&eventModel.UserUUID,
		&eventModel.Title,
		&eventModel.Description,
		&eventModel.StartDate,
		&eventModel.EndDate,
		&eventModel.DelayNotification,
		&eventModel.DelayNotificationType,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEventNotFoundByGet
	}

	return &eventModel, err
}

func (repo *repositoryDB) List(ctx context.Context, filter Filter) ([]*event.Event, error) {
	err := CheckFilter(filter)
	if err != nil {
		return nil, err
	}

	query, args := buildQuery(filter, false)
	rows, err := repo.db.QueryContext(ctx, query, args...)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if err != nil {
		return nil, err
	}

	var events []*event.Event
	for rows.Next() {
		var eventModel event.Event
		err = rows.Scan(
			&eventModel.UUID,
			&eventModel.UserUUID,
			&eventModel.Title,
			&eventModel.Description,
			&eventModel.StartDate,
			&eventModel.EndDate,
			&eventModel.DelayNotification,
			&eventModel.DelayNotificationType,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &eventModel)
	}
	return events, rows.Err()
}

func (repo *repositoryDB) Exists(ctx context.Context, filter Filter) (bool, error) {
	query, args := buildQuery(filter, false) // без LIMIT
	query = strings.Replace(query, "SELECT *", "SELECT 1", 1)

	query = fmt.Sprintf("SELECT EXISTS (%s)", query)

	var exists bool
	err := repo.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func buildQuery(filter Filter, isOne bool) (string, []interface{}) {
	query := `SELECT * FROM event `

	var where []string
	var args []interface{}

	argIndex := 1
	if filter.UUID != nil {
		where = append(where, fmt.Sprintf("uuid = $%d", argIndex))
		args = append(args, *filter.UUID)
		argIndex++
	}
	if filter.UserUUID != nil {
		where = append(where, fmt.Sprintf("user_uuid = $%d", argIndex))
		args = append(args, *filter.UserUUID)
		argIndex++
	}

	var busyParts []string
	if filter.StartDateByBusy != nil {
		busyParts = append(busyParts, fmt.Sprintf("(start_date <= $%d AND end_date >= $%d)", argIndex, argIndex))
		args = append(args, *filter.StartDateByBusy)
		argIndex++
	}
	if filter.EndDateByBusy != nil {
		busyParts = append(busyParts, fmt.Sprintf("(start_date <= $%d AND end_date >= $%d)", argIndex, argIndex))
		args = append(args, *filter.EndDateByBusy)
		argIndex++ //nolint:ineffassign
	}
	if len(busyParts) > 0 {
		where = append(where, "("+strings.Join(busyParts, " OR ")+")")
	}

	if len(where) > 0 {
		query += "WHERE " + strings.Join(where, " AND ") + "\n"
	}

	if isOne {
		query += "LIMIT 1"
	}

	return query, args
}
