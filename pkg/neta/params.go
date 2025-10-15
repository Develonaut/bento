package neta

// GetStringParam extracts a string parameter with default.
func GetStringParam(params map[string]interface{}, key, defaultVal string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultVal
}
