package usecase

import (
	"fmt"
	"log"
	"time"

	"Name_IQ_Finder/internal/entity"
)

type PersonUseCase struct {
	repo   entity.PersonRepository
	client entity.ExternalAPIClient
	logger *log.Logger
}

func NewPersonUseCase(repo entity.PersonRepository, client entity.ExternalAPIClient, logger *log.Logger) *PersonUseCase {
	return &PersonUseCase{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

func (uc *PersonUseCase) Create(name, surname, patronymic string) (*entity.Person, error) {
	uc.logger.Printf("Creating person with name=%s, surname=%s", name, surname)

	age, gender, nationality, err := uc.client.EnrichPerson(name)
	if err != nil {
		uc.logger.Printf("Error enriching person data: %v", err)
		return nil, fmt.Errorf("failed to enrich person data: %w", err)
	}

	uc.logger.Printf("Enriched data: age=%d, gender=%s, nationality=%s", age, gender, nationality)

	person := &entity.Person{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	id, err := uc.repo.Create(person)
	if err != nil {
		uc.logger.Printf("Error creating person in repository: %v", err)
		return nil, fmt.Errorf("failed to create person: %w", err)
	}

	person.ID = id
	uc.logger.Printf("Person created with ID=%d", id)

	return person, nil
}

func (uc *PersonUseCase) GetByID(id int64) (*entity.Person, error) {
	uc.logger.Printf("Getting person with ID=%d", id)

	person, err := uc.repo.GetByID(id)
	if err != nil {
		uc.logger.Printf("Error getting person with ID=%d: %v", id, err)
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	uc.logger.Printf("Found person with ID=%d", id)
	return person, nil
}

func (uc *PersonUseCase) GetAll(filter map[string]interface{}, page, limit int) ([]*entity.Person, int, error) {
	uc.logger.Printf("Getting persons with filter=%v, page=%d, limit=%d", filter, page, limit)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	persons, total, err := uc.repo.GetAll(filter, page, limit)
	if err != nil {
		uc.logger.Printf("Error getting persons: %v", err)
		return nil, 0, fmt.Errorf("failed to get persons: %w", err)
	}

	uc.logger.Printf("Found %d persons (total: %d)", len(persons), total)
	return persons, total, nil
}

func (uc *PersonUseCase) Update(person *entity.Person) error {
	uc.logger.Printf("Updating person with ID=%d", person.ID)

	err := uc.repo.Update(person)
	if err != nil {
		uc.logger.Printf("Error updating person with ID=%d: %v", person.ID, err)
		return fmt.Errorf("failed to update person: %w", err)
	}

	uc.logger.Printf("Person with ID=%d updated", person.ID)
	return nil
}

func (uc *PersonUseCase) Delete(id int64) error {
	uc.logger.Printf("Deleting person with ID=%d", id)

	err := uc.repo.Delete(id)
	if err != nil {
		uc.logger.Printf("Error deleting person with ID=%d: %v", id, err)
		return fmt.Errorf("failed to delete person: %w", err)
	}

	uc.logger.Printf("Person with ID=%d deleted", id)
	return nil
}
