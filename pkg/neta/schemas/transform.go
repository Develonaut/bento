package schemas

// JQSchema validates jq transform node parameters.
type JQSchema struct{}

// NewJQSchema creates a jq schema validator.
func NewJQSchema() *JQSchema {
	return &JQSchema{}
}

// Validate checks jq transform parameters.
func (s *JQSchema) Validate(params map[string]interface{}) error {
	var errs ValidationErrors

	// Required: query
	query, ok := params["query"].(string)
	if !ok || query == "" {
		errs = append(errs, ValidationError{
			Field:   "query",
			Message: "is required (jq query string)",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Fields returns field definitions for form generation.
func (s *JQSchema) Fields() []Field {
	return []Field{
		{
			Name:        "query",
			Type:        FieldString,
			Required:    true,
			Description: "jq query to transform data (e.g., '.data | .[] | select(.active)')",
		},
		{
			Name:        "input",
			Type:        FieldString,
			Required:    false,
			Description: "Static input data (if not using previous node output)",
		},
	}
}
