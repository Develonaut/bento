package schemas

// FileWriterSchema validates file.write node parameters.
type FileWriterSchema struct{}

// NewFileWriterSchema creates a file writer schema validator.
func NewFileWriterSchema() *FileWriterSchema {
	return &FileWriterSchema{}
}

// Validate checks file writer parameters.
func (s *FileWriterSchema) Validate(params map[string]interface{}) error {
	var errs ValidationErrors

	// Required: path
	path, ok := params["path"].(string)
	if !ok || path == "" {
		errs = append(errs, ValidationError{
			Field:   "path",
			Message: "is required (file path to write to)",
		})
	}

	// Content is optional - can come from previous node's input
	// No validation needed for content

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Fields returns field definitions for form generation.
func (s *FileWriterSchema) Fields() []Field {
	return []Field{
		{
			Name:        "path",
			Type:        FieldString,
			Required:    true,
			Description: "File path to write to (e.g., '/tmp/output.txt' or '~/myfile.json')",
		},
		{
			Name:        "content",
			Type:        FieldString,
			Required:    false,
			Description: "Content to write to the file (optional, uses previous node output if not specified)",
		},
	}
}
