package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type user struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

var _ UserStore = (*user)(nil)

func NewUser(storage *PGres, logger *zap.SugaredLogger) UserStore {
	return &user{
		storage: storage,
		logger:  logger,
	}
}

func (u *user) Create(ctx context.Context, user *models.User) (*models.User, error) {
	u.logger.Debugf("Enter in user repository Create() with args: ctx, user: %v", user)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		// we create rights and address somewhere in usecase or get them from user
		row := pool.QueryRow(ctx, `INSERT INTO users 
		(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
			user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
			user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)
		var id uuid.UUID
		err := row.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("can't create user: %w", err)
		}
		user.ID = id
		return user, nil
	}
}

func (u *user) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u.logger.Debug("Enter in repository GetUserByEmail() with args: ctx, email: %v", email)
	select {
	case <-ctx.Done():
		return &models.User{}, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		row := pool.QueryRow(ctx, `SELECT users.id, users.name, lastname, password, email, rights.id, zipcode, country, city, street, rights.name, rights.rules FROM users INNER JOIN rights ON email=$1 and rights.id=users.rights`, email)
		var user = models.User{}
		err := row.Scan(
			&user.ID,
			&user.Firstname,
			&user.Lastname,
			&user.Password,
			&user.Email,
			&user.Rights.ID,
			&user.Address.Zipcode,
			&user.Address.Country,
			&user.Address.City,
			&user.Address.Street,
			&user.Rights.Name,
			&user.Rights.Rules)
		if err != nil {
			return &models.User{}, fmt.Errorf("can't get user from database: %w", err)
		}
		u.logger.Info("Get user by email success")
		return &user, nil
	}
}

func (u *user) UpdateUserData(ctx context.Context, user *models.User) (*models.User, error) {
	u.logger.Debug("Enter in repository UpdateUserData()")
	select {
	case <-ctx.Done():
		return &models.User{}, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			u.logger.Errorf("can't create transaction: %s", err)
			return &models.User{}, fmt.Errorf("can't create transaction: %w", err)
		}
		u.logger.Debug("transaction begin success")
		defer func() {
			if err != nil {
				u.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					u.logger.Errorf("can't rollback %s", err)
				}

			} else {
				u.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					u.logger.Errorf("can't commit %s", err)
				}
			}
		}()

		_, err = tx.Exec(ctx, `UPDATE users SET name=$1, lastname=$2, country=$3, city=$4, street=$5, zipcode=$6 WHERE id=$7`,
			user.Firstname,
			user.Lastname,
			user.Address.Country,
			user.Address.City,
			user.Address.Street,
			user.Address.Zipcode,
			user.ID)
		if err != nil {
			u.logger.Errorf("error on update user %s: %s", user.ID, err)
			return &models.User{}, fmt.Errorf("error on update user %s: %w", user.ID, err)
		}
		u.logger.Infof("user %s successfully updated %s", user.ID, user.Lastname)
		return user, nil
	}
}

func (u *user) GetRightsId(ctx context.Context, name string) (models.Rights, error) {
	select {
	case <-ctx.Done():
		return models.Rights{}, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		row := pool.QueryRow(ctx, `SELECT id, name, rules FROM rights WHERE name=$1`, name)
		var rights = models.Rights{}
		err := row.Scan(&rights.ID, &rights.Name, &rights.Rules)
		if err != nil {
			return models.Rights{}, fmt.Errorf("can't get rights from database: %w", err)
		}
		return rights, nil

	}
}

func (u *user) SaveSession(ctx context.Context, token string, t int64) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		pool.QueryRow(ctx, `INSERT INTO session (token, timestamp) VALUES ($1, $2)`,
			token, t)
	}
	return nil
}

func (u *user) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u.logger.Debugf("Enter in repository GetUserById() with args: ctx, id: %v", id)
	select {
	case <-ctx.Done():
		return &models.User{}, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		row := pool.QueryRow(ctx, `SELECT users.id, users.name, lastname, password, email, rights.id, zipcode, country, city, street, rights.name, rights.rules FROM users INNER JOIN rights ON users.id=$1 and rights.id=users.rights`, id)
		var user = models.User{}
		err := row.Scan(
			&user.ID,
			&user.Firstname,
			&user.Lastname,
			&user.Password,
			&user.Email,
			&user.Rights.ID,
			&user.Address.Zipcode,
			&user.Address.Country,
			&user.Address.City,
			&user.Address.Street,
			&user.Rights.Name,
			&user.Rights.Rules)
		if err != nil {
			return &models.User{}, fmt.Errorf("can't get user from database: %w", err)
		}
		return &user, nil
	}
}

func (u *user) GetUsersList(ctx context.Context) ([]models.User, error) {
	u.logger.Debug("Enter in repository GetAllUsers()")

	pool := u.storage.GetPool()
	user := models.User{}
	usersList := make([]models.User, 0, 100)

	rows, err := pool.Query(ctx, `SELECT users.id, users.name, lastname, email, rights.id, zipcode, country, city, street, rights.name, rights.rules FROM users INNER JOIN rights ON rights.id=users.rights`)
	if err != nil {
		msg := fmt.Errorf("error on users list query context: %w", err)
		u.logger.Error(msg.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&user.ID,
			&user.Firstname,
			&user.Lastname,
			&user.Email,
			&user.Rights.ID,
			&user.Address.Zipcode,
			&user.Address.Country,
			&user.Address.City,
			&user.Address.Street,
			&user.Rights.Name,
			&user.Rights.Rules); err != nil {
			u.logger.Error(err.Error())
			return nil, err
		}
		usersList = append(usersList, user)
	}
	return usersList, nil
}

func (u *user) ChangeUserRole(ctx context.Context, userId uuid.UUID, rightsId uuid.UUID) error {
	u.logger.Debug("Enter in repository ChangeUserRole() with args: ctx, userId: %v, rightsId: %v", userId, rightsId)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			u.logger.Errorf("can't create transaction: %s", err)
			return fmt.Errorf("can't create transaction: %w", err)
		}
		u.logger.Debug("transaction begin success")
		defer func() {
			if err != nil {
				u.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					u.logger.Errorf("can't rollback %s", err)
				}

			} else {
				u.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					u.logger.Errorf("can't commit %s", err)
				}
			}
		}()

		_, err = tx.Exec(ctx, `UPDATE users SET rights=$1 WHERE id=$2`,
			rightsId,
			userId)
		if err != nil {
			u.logger.Errorf("error on update user %s: %s", userId, err)
			return fmt.Errorf("error on update user %s: %w", userId, err)
		}
		u.logger.Infof("user with id %s successfully updated: new rights is: %s", userId, rightsId)
		return nil
	}
}

func (u *user) ChangeUserPassword(ctx context.Context, userId uuid.UUID, newPassword string) error {
	u.logger.Debug("Enter in repository ChangeUserPassword() with args: ctx, userId: %v, newPassword: %s", userId, newPassword)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			u.logger.Errorf("can't create transaction: %s", err)
			return fmt.Errorf("can't create transaction: %w", err)
		}
		u.logger.Debug("transaction begin success")
		defer func() {
			if err != nil {
				u.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					u.logger.Errorf("can't rollback %s", err)
				}

			} else {
				u.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					u.logger.Errorf("can't commit %s", err)
				}
			}
		}()

		_, err = tx.Exec(ctx, `UPDATE users SET password=$1 WHERE id=$2`,
			newPassword,
			userId)
		if err != nil {
			u.logger.Errorf("error on update user %s: %s", userId, err)
			return fmt.Errorf("error on update user %s: %w", userId, err)
		}
		u.logger.Infof("user with id %s successfully updated", userId)
		return nil
	}
}

func (u *user) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	u.logger.Debug("Enter in repository DeleteUser() with args: ctx, userId: %v", userId)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := u.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		defer func() {
			if err != nil {
				u.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					u.logger.Errorf("can't rollback %s", err)
				}

			} else {
				u.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					u.logger.Errorf("can't commit %s", err)
				}
			}
		}()
		_, err = tx.Exec(ctx, `DELETE FROM users WHERE id=$1`, userId)
		if err != nil {
			u.logger.Errorf("can't delete user: %s", err)
			return fmt.Errorf("can't delete user: %w", err)
		}
		u.logger.Info("Delete user with id: %v from database success", userId)
		return nil
	}
}
