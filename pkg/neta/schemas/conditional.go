package schemas

// IfSchema validates conditional.if node parameters.
type IfSchema struct{}

// NewIfSchema creates a conditional.if schema validator.
func NewIfSchema() *IfSchema {
	return &IfSchema{}
}

// Validate checks conditional.if parameters.
func (s *IfSchema) Validate(params map[string]interface{}) error {
	var errs ValidationErrors

	// Required: condition (must be bool)
	if _, ok := params["condition"].(bool); !ok {
		errs = append(errs, ValidationError{
			Field:   "condition",
			Message: "is required and must be a boolean",
		})
	}

	// Optional: then (should be a Definition, but we can't validate structure here)
	// Optional: else (should be a Definition, but we can't validate structure here)
	// At least one branch should be present, but both are technically optional

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Fields returns field definitions for form generation.
func (s *IfSchema) Fields() []Field {
	return []Field{
		{
			Name:        "condition",
			Type:        FieldBool,
			Required:    true,
			Description: "Boolean condition to evaluate",
		},
		{
			Name:        "then",
			Type:        FieldMap,
			Required:    false,
			Description: "Node definition to execute if condition is true",
		},
		{
			Name:        "else",
			Type:        FieldMap,
			Required:    false,
			Description: "Node definition to execute if condition is false",
		},
	}
}
