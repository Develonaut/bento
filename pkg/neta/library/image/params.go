package image

// getIntParam extracts an integer parameter with a default value.
func getIntParam(params map[string]interface{}, key string, defaultVal int) int {
	if val, ok := params[key].(int); ok {
		return val
	}
	if val, ok := params[key].(float64); ok {
		return int(val)
	}
	return defaultVal
}

// getBoolParam extracts a boolean parameter with a default value.
func getBoolParam(params map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := params[key].(bool); ok {
		return val
	}
	return defaultVal
}

// getStringParam extracts a string parameter with a default value.
func getStringParam(params map[string]interface{}, key string, defaultVal string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultVal
}
