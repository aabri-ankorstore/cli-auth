package repository

type Repository interface {
	Get(ID string) (interface{}, error)
	GetAll() ([]interface{}, error)
	Insert(value interface{}) error
	Update(ID string, data interface{}, query string) error
	Delete(ID string) error
}
