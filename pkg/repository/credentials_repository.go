package repository

import (
	"context"
	"github.com/aabri-ankorstore/cli-auth/pkg/entities"
	"github.com/uptrace/bun"
)

type CredentialsRepository struct {
	DB  bun.IDB
	Ctx context.Context
}

func (c *CredentialsRepository) Get(ID string) (interface{}, error) {
	var credential entities.Credential
	if err := c.DB.NewSelect().Model(&credential).Where("account_id = ?", ID).Scan(c.Ctx); err != nil {
		return nil, err
	}
	return credential, nil
}
func (c *CredentialsRepository) GetAll() ([]interface{}, error) {
	var credentials []entities.Credential
	if err := c.DB.NewSelect().Model(&credentials).Scan(c.Ctx); err != nil {
		return nil, err
	}
	var data []interface{}
	for _, v := range credentials {
		data = append(data, v)
	}
	return data, nil
}
func (c *CredentialsRepository) Insert(value interface{}) error {
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
func (c *CredentialsRepository) Update(value interface{}) error {
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
func (c *CredentialsRepository) Delete(ID string) error {
	credential := new(entities.Credential)
	_, err := c.DB.NewDelete().Model(credential).Where("account_id = ?", ID).Exec(c.Ctx)
	return err
}
