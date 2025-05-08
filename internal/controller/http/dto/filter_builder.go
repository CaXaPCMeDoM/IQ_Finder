package dto

type FilterBuilder struct {
	filter map[string]interface{}
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		filter: make(map[string]interface{}),
	}
}

func (b *FilterBuilder) WithName(name string) *FilterBuilder {
	if name != "" {
		b.filter["name"] = name
	}
	return b
}

func (b *FilterBuilder) WithSurname(surname string) *FilterBuilder {
	if surname != "" {
		b.filter["surname"] = surname
	}
	return b
}

func (b *FilterBuilder) WithNationality(nationality string) *FilterBuilder {
	if nationality != "" {
		b.filter["nationality"] = nationality
	}
	return b
}

func (b *FilterBuilder) Build() map[string]interface{} {
	return b.filter
}
