package entity

type PersonRepository interface {
	Create(person *Person) (int64, error)
	GetByID(id int64) (*Person, error)
	GetAll(filter map[string]interface{}, page, limit int) ([]*Person, int, error)
	Update(person *Person) error
	Delete(id int64) error
}
