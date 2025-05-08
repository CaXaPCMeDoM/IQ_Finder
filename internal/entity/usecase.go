package entity

type PersonUseCase interface {
	Create(name, surname, patronymic string) (*Person, error)
	GetByID(id int64) (*Person, error)
	GetAll(filter map[string]interface{}, page, limit int) ([]*Person, int, error)
	Update(person *Person) error
	Delete(id int64) error
}

type ExternalAPIClient interface {
	GetAge(name string) (int, error)
	GetGender(name string) (string, error)
	GetNationality(name string) (string, error)
	EnrichPerson(name string) (int, string, string, error)
}
