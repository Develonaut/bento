package schemas

// SequenceSchema validates sequence group parameters.
type SequenceSchema struct{}

// NewSequenceSchema creates a sequence schema validator.
func NewSequenceSchema() *SequenceSchema {
	return &SequenceSchema{}
}

// Validate checks sequence parameters.
// Sequences have no parameters - validation is structural.
func (s *SequenceSchema) Validate(params map[string]interface{}) error {
	return nil
}

// Fields returns field definitions for form generation.
func (s *SequenceSchema) Fields() []Field {
	return []Field{
		{
			Name:        "nodes",
			Type:        FieldArray,
			Required:    true,
			Description: "Child nodes to execute in sequence",
		},
	}
}

// ParallelSchema validates parallel group parameters.
type ParallelSchema struct{}

// NewParallelSchema creates a parallel schema validator.
func NewParallelSchema() *ParallelSchema {
	return &ParallelSchema{}
}

// Validate checks parallel parameters.
func (s *ParallelSchema) Validate(params map[string]interface{}) error {
	var errs ValidationErrors

	// Optional: max_concurrent (must be non-negative int)
	if maxConcurrent, ok := params["max_concurrent"]; ok {
		if err := validateMaxConcurrent(maxConcurrent); err != nil {
			errs = append(errs, ValidationError{
				Field:   "max_concurrent",
				Message: err.Error(),
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Fields returns field definitions for form generation.
func (s *ParallelSchema) Fields() []Field {
	defaultZero := 0
	return []Field{
		{
			Name:        "nodes",
			Type:        FieldArray,
			Required:    true,
			Description: "Child nodes to execute in parallel",
		},
		{
			Name:        "max_concurrent",
			Type:        FieldInt,
			Required:    false,
			Description: "Maximum concurrent executions (0 = unlimited)",
			Default:     defaultZero,
			Min:         &defaultZero,
		},
	}
}

// validateMaxConcurrent checks if value is a valid max_concurrent.
func validateMaxConcurrent(val interface{}) error {
	intVal, ok := val.(int)
	if !ok {
		return ValidationError{
			Field:   "max_concurrent",
			Message: "must be an integer",
		}
	}

	if intVal < 0 {
		return ValidationError{
			Field:   "max_concurrent",
			Message: "must be non-negative",
		}
	}

	return nil
}
