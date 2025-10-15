// Package transform provides data transformation nodes.
package transform

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"

	"bento/pkg/neta"
)

// JQ applies jq transformations to data.
type JQ struct{}

// NewJQ creates a new JQ transformer.
func NewJQ() *JQ {
	return &JQ{}
}

// Execute applies a jq query to input data.
func (j *JQ) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	query := neta.GetStringParam(params, "query", ".")
	input := params["input"]

	// If input is a string, try to parse as JSON
	if strInput, ok := input.(string); ok {
		var parsed interface{}
		if err := json.Unmarshal([]byte(strInput), &parsed); err != nil {
			return neta.Result{}, fmt.Errorf("failed to parse input JSON: %w", err)
		}
		input = parsed
	}

	result, err := applyQuery(query, input)
	if err != nil {
		return neta.Result{}, fmt.Errorf("jq transform failed: %w", err)
	}

	return neta.Result{Output: result}, nil
}

// applyQuery executes a jq query on data.
func applyQuery(queryStr string, data interface{}) (interface{}, error) {
	query, err := gojq.Parse(queryStr)
	if err != nil {
		return nil, err
	}

	iter := query.Run(data)
	v, ok := iter.Next()
	if !ok {
		return nil, fmt.Errorf("no result from query")
	}
	if err, ok := v.(error); ok {
		return nil, err
	}

	return v, nil
}
