package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

var _ RightsStore = rightsRepo{}

type rightsRepo struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

func NewRightsRepo(storage *PGres, logger *zap.SugaredLogger) RightsStore {
	return &rightsRepo{
		storage: storage,
		logger:  logger,
	}
}

func (repo rightsRepo) CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error) {
	repo.logger.Debugf("Enter in repository CreateRights() with args: ctx, rights: %v", rights)

	var id uuid.UUID
	pool := repo.storage.GetPool()

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("Can't create transaction: %s", err)
		return uuid.Nil, fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("Can't rollback %s", err)
			}

		} else {
			repo.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("Can't commit %s", err)
			}
		}
	}()
	row := tx.QueryRow(ctx, `INSERT INTO rights(name, rules) values ($1, $2) RETURNING id`,
		rights.Name,
		rights.Rules,
	)
	err = row.Scan(&id)
	if err != nil {
		repo.logger.Errorf("can't create rights %s", err)
		return uuid.Nil, fmt.Errorf("can't create rights %w", err)
	}
	repo.logger.Info("Rights create success")
	repo.logger.Debugf("id is %v\n", id)
	return id, nil
}

func (repo rightsRepo) UpdateRights(ctx context.Context, rights *models.Rights) error {
	repo.logger.Debugf("Enter in repo UpdateRights() with args: ctx, rights: %v", rights)

	pool := repo.storage.GetPool()

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("Can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("Can't rollback %s", err)
			}
		} else {
			repo.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("Can't commit %s", err)
			}
		}
	}()

	_, err = tx.Exec(ctx, `UPDATE rights SET name=$1, rules=$2 WHERE id=$3`,
		rights.Name,
		rights.Rules,
		rights.ID)
	if err != nil {
		repo.logger.Errorf("Error on update rights %s: %s", rights.ID, err)
		return fmt.Errorf("error on update rights %s: %w", rights.ID, err)
	}
	repo.logger.Infof("Item %s successfully updated", rights.ID)
	return nil
}

func (repo rightsRepo) DeleteRights(ctx context.Context, id uuid.UUID) error {
	repo.logger.Debugf("Enter in repo DeleteRights() with args: ctx, id: %v", id)

	pool := repo.storage.GetPool()
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			repo.logger.Errorf("transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("can't rollback %s", err)
			}

		} else {
			repo.logger.Info("transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("can't commit %s", err)
			}
		}
	}()
	_, err = tx.Exec(ctx, `DELETE FROM rights WHERE id=$1`, id)
	if err != nil {
		repo.logger.Errorf("can't delete rights: %s", err)
		return fmt.Errorf("can't delete cart rights: %w", err)
	}
	repo.logger.Info("Delete rights with id: %v from database success", id)
	return nil

}

func (repo rightsRepo) GetRights(ctx context.Context, id uuid.UUID) (*models.Rights, error) {
	repo.logger.Debugf("Enter in repository GetRights() with args ctx, id: %v", id)

	rights := models.Rights{}
	pool := repo.storage.GetPool()
	row := pool.QueryRow(ctx, `SELECT id, name, rules FROM rights WHERE id=$1`, id)
	err := row.Scan(
		&rights.ID,
		&rights.Name,
		&rights.Rules,
	)
	if err != nil {
		repo.logger.Errorf("Error in rows scan get rights by id: %s", err)
		return &models.Rights{}, fmt.Errorf("error in rows scan get rights by id: %w", err)
	}
	repo.logger.Info("Get rights success")
	return &rights, nil
}

func (repo rightsRepo) RightsList(ctx context.Context) ([]models.Rights, error) {
	repo.logger.Debugf("Enter in repository RightsList() with args: ctx")

	pool := repo.storage.GetPool()
	rights := models.Rights{}
	rightsList := make([]models.Rights, 0, 100)

	rows, err := pool.Query(ctx, `SELECT * FROM rights`)
	if err != nil {
		msg := fmt.Errorf("error on rights list query context: %w", err)
		repo.logger.Error(msg.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&rights.ID,
			&rights.Name,
			&rights.Rules,
		); err != nil {
			repo.logger.Error(err.Error())
			return nil, err
		}
		rightsList = append(rightsList, rights)
	}
	return rightsList, nil
}

func (repo rightsRepo) GetRightsByName(ctx context.Context, name string) (*models.Rights, error) {
	repo.logger.Debugf("Enter in repository GetRightsByName() with args ctx, name: %s", name)

	rights := models.Rights{}
	pool := repo.storage.GetPool()
	row := pool.QueryRow(ctx, `SELECT id, name, rules FROM rights WHERE name=$1`, name)
	err := row.Scan(
		&rights.ID,
		&rights.Name,
		&rights.Rules,
	)
	if err != nil {
		repo.logger.Errorf("Error in rows scan get rights by name: %s", err)
		return &models.Rights{}, fmt.Errorf("error in rows scan get rights by name: %w", err)
	}
	repo.logger.Info("Get rights success")
	return &rights, nil
}
