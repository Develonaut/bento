package schemas

// ForLoopSchema validates loop.for node parameters.
type ForLoopSchema struct{}

// NewForLoopSchema creates a for loop schema validator.
func NewForLoopSchema() *ForLoopSchema {
	return &ForLoopSchema{}
}

// Validate checks loop.for parameters.
func (s *ForLoopSchema) Validate(params map[string]interface{}) error {
	var errs ValidationErrors

	// Required: items (must be array)
	// Note: Empty arrays are valid - they result in no iterations
	if _, ok := params["items"].([]interface{}); !ok {
		errs = append(errs, ValidationError{
			Field:   "items",
			Message: "is required and must be an array",
		})
	}

	// Required: body (must be a Definition)
	// Note: In YAML this will be parsed as a map, but at execution time
	// it's converted to a Definition. We can't fully validate structure here.
	if _, ok := params["body"]; !ok {
		errs = append(errs, ValidationError{
			Field:   "body",
			Message: "is required (node definition to execute for each item)",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Fields returns field definitions for form generation.
func (s *ForLoopSchema) Fields() []Field {
	return []Field{
		{
			Name:        "items",
			Type:        FieldArray,
			Required:    true,
			Description: "Array of items to iterate over",
		},
		{
			Name:        "body",
			Type:        FieldMap,
			Required:    true,
			Description: "Node definition to execute for each item",
		},
	}
}
