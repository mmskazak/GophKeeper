package pwd_services

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gophKeeper/internal/modules/pwd/pwd_dto"
)

type PwdService struct {
	pool *pgxpool.Pool
}

func NewPwdService(pool *pgxpool.Pool) *PwdService {
	return &PwdService{pool: pool}
}

func (pwd *PwdService) SavePassword(ctx context.Context, dto pwd_dto.SavePwdDTO) error {
	sql := `INSERT INTO passwords (user_id, resource, login, password) VALUES ($1, $2, $3, $4)`
	_, err := pwd.pool.Exec(ctx, sql, dto.UserID)
	if err != nil {
		return fmt.Errorf("error save password from pwd service: %w", err)
	}

	return nil
}

func (pwd *PwdService) DeletePassword(ctx context.Context, dto pwd_dto.DeletePwdDTO) error {
	sql := `DELETE FROM passwords WHERE user_id = $1 AND resource = $2;`
	_, err := pwd.pool.Exec(ctx, sql, dto.UserID)
	if err != nil {
		return fmt.Errorf("error save password from pwd service: %w", err)
	}
	return nil
}

func (pwd *PwdService) GetPassword(ctx context.Context, dto pwd_dto.GetPwdDTO) (string, error) {
	sql := `SELECT password FROM passwords WHERE user_id = $1 AND id = $2;`
	row := pwd.pool.QueryRow(ctx, sql, dto.UserID, dto.PwdID)

	var password string
	err := row.Scan(&password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Обработка случая, когда запись не найдена
			return "", fmt.Errorf(
				"password not found for user_id %v and pwd_id %v: %w",
				dto.UserID,
				dto.PwdID,
				err,
			)
		}
		// Обработка других ошибок
		return "", fmt.Errorf("error scanning password from pwd service: %w", err)
	}

	return password, nil
}

func (pwd *PwdService) GetAllPasswords(ctx context.Context, dto pwd_dto.AllPwdDTO) ([]InfoByPassword, error) {
	sql := `SELECT resource, login, password FROM passwords WHERE user_id = $1`
	rows, err := pwd.pool.Query(ctx, sql, dto.UserID)
	if err != nil {
		return []InfoByPassword{}, fmt.Errorf("error query get all passwords: %w", err)
	}

	var listPasswords []InfoByPassword

	for rows.Next() {
		var resource string
		var login string
		var password string
		err := rows.Scan(&resource, &login, &password)
		if err != nil {
			return []InfoByPassword{}, fmt.Errorf("error scan get all passwords from pwd service: %w", err)
		}

		listPasswords = append(listPasswords, InfoByPassword{
			Resource: resource,
			Login:    login,
			Password: password,
		})
	}

	return listPasswords, nil
}
