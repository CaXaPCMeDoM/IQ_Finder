package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Name_IQ_Finder/internal/controller/http/dto"
	"Name_IQ_Finder/internal/entity"
)

type PersonHandler struct {
	useCase entity.PersonUseCase
}

func NewPersonHandler(useCase entity.PersonUseCase) *PersonHandler {
	return &PersonHandler{
		useCase: useCase,
	}
}

// Create @Summary Create person
// @Description Create a new person with enriched data
// @Tags persons
// @Accept json
// @Produce json
// @Param request body dto.CreatePersonRequest true "Person data"
// @Success 201 {object} dto.PersonResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/persons [post]
func (h *PersonHandler) Create(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := h.useCase.Create(req.Name, req.Surname, req.Patronymic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toPersonResponse(person))
}

// GetByID @Summary Get person by ID
// @Description Get a person by ID
// @Tags persons
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} dto.PersonResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/persons/{id} [get]
func (h *PersonHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	person, err := h.useCase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toPersonResponse(person))
}

// GetAll @Summary Get all persons
// @Description Get all persons with optional filtering and pagination
// @Tags persons
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param name query string false "Filter by name"
// @Param surname query string false "Filter by surname"
// @Param nationality query string false "Filter by nationality"
// @Success 200 {object} dto.PersonListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/persons [get]
func (h *PersonHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := dto.NewFilterBuilder().
		WithName(c.Query("name")).
		WithSurname(c.Query("surname")).
		WithNationality(c.Query("nationality")).
		Build()

	persons, total, err := h.useCase.GetAll(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	personResponses := make([]dto.PersonResponse, len(persons))
	for i, person := range persons {
		personResponses[i] = toPersonResponse(person)
	}

	c.JSON(http.StatusOK, dto.PersonListResponse{
		Data:       personResponses,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	})
}

// Update @Summary Update person
// @Description Update an existing person
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param request body dto.UpdatePersonRequest true "Person data"
// @Success 200 {object} dto.PersonResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/persons/{id} [put]
func (h *PersonHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	person, err := h.useCase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		person.Name = req.Name
	}
	if req.Surname != "" {
		person.Surname = req.Surname
	}
	if req.Patronymic != "" {
		person.Patronymic = req.Patronymic
	}
	if req.Age != nil {
		person.Age = *req.Age
	}
	if req.Gender != "" {
		person.Gender = req.Gender
	}
	if req.Nationality != "" {
		person.Nationality = req.Nationality
	}

	if err := h.useCase.Update(person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toPersonResponse(person))
}

// Delete @Summary Delete person
// @Description Delete a person by ID
// @Tags persons
// @Param id path int true "Person ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/persons/{id} [delete]
func (h *PersonHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.useCase.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func toPersonResponse(person *entity.Person) dto.PersonResponse {
	return dto.PersonResponse{
		ID:          person.ID,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
		CreatedAt:   person.CreatedAt,
		UpdatedAt:   person.UpdatedAt,
	}
}
