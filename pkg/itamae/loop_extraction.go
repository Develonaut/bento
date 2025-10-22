package itamae

import (
	"fmt"

	"github.com/Develonaut/bento/pkg/neta"
)

// extractLoopItems extracts and validates items for forEach loop.
func (i *Itamae) extractLoopItems(
	def *neta.Definition,
	execCtx *executionContext,
) ([]interface{}, error) {
	itemsParam := def.Parameters["items"]

	if i.logger != nil {
		i.logger.Debug("Loop items parameter",
			"loop_id", def.ID,
			"itemsParam", itemsParam,
			"itemsParam_type", fmt.Sprintf("%T", itemsParam))
	}

	resolved := execCtx.resolveValue(itemsParam)

	if i.logger != nil {
		i.logger.Debug("Loop items resolved",
			"loop_id", def.ID,
			"resolved", resolved,
			"resolved_type", fmt.Sprintf("%T", resolved))
	}

	return i.convertToInterfaceArray(def, resolved)
}

// convertToInterfaceArray converts resolved value to []interface{}.
func (i *Itamae) convertToInterfaceArray(
	def *neta.Definition,
	resolved interface{},
) ([]interface{}, error) {
	switch v := resolved.(type) {
	case []interface{}:
		return v, nil
	case []map[string]interface{}:
		items := make([]interface{}, len(v))
		for idx, item := range v {
			items[idx] = item
		}
		return items, nil
	default:
		if i.logger != nil {
			i.logger.Error("Loop items not an array",
				"loop_id", def.ID,
				"resolved_type", fmt.Sprintf("%T", resolved),
				"resolved_value", resolved)
		}
		return nil, newNodeError(def.ID, "loop", "validate",
			fmt.Errorf("'items' must be an array, got %T", resolved))
	}
}
