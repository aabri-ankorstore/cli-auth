package repository

import (
	"context"
	"github.com/aabri-ankorstore/cli-auth/pkg/entities"
	"github.com/uptrace/bun"
)

type AccessTokensRepository struct {
	DB  bun.IDB
	Ctx context.Context
}

func (c *AccessTokensRepository) Get(ID string) (interface{}, error) {
	var accessToken entities.AccessToken
	if err := c.DB.NewSelect().Model(&accessToken).Where("account_id = ?", ID).Scan(c.Ctx); err != nil {
		return nil, err
	}
	return accessToken, nil
}
func (c *AccessTokensRepository) GetAll() ([]interface{}, error) {
	var accessToken []entities.AccessToken
	if err := c.DB.NewSelect().Model(&accessToken).Scan(c.Ctx); err != nil {
		return nil, err
	}
	var data []interface{}
	for _, v := range accessToken {
		data = append(data, v)
	}
	return data, nil
}
func (c *AccessTokensRepository) Insert(value interface{}) error {
	switch values := value.(type) {
	case []interface{}:
		for v := range values {
			if _, err := c.DB.NewInsert().
				Model(v).
				Exec(c.Ctx); err != nil {
				return err
			}
		}
	case interface{}:
		_, err := c.DB.NewInsert().
			Model(value).
			Exec(c.Ctx)
		return err
	}
	return nil
}
func (c *AccessTokensRepository) Update(value interface{}) error {
	switch values := value.(type) {
	case []interface{}:
		for v := range values {
			if _, err := c.DB.NewUpdate().
				Model(v).
				WherePK().
				Exec(c.Ctx); err != nil {
				return err
			}
		}
	case interface{}:
		_, err := c.DB.NewUpdate().
			Model(value).
			WherePK().
			Exec(c.Ctx)
		return err
	}
	return nil
}
func (c *AccessTokensRepository) Delete(ID string) error {
	accessToken := new(entities.AccessToken)
	_, err := c.DB.NewDelete().Model(accessToken).Where("account_id = ?", ID).Exec(c.Ctx)
	return err
}
