package template

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// replace performs template replacement operations.
func (t *Template) replace(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	input, output, replMap, mode, err := validateAndExtractParams(params)
	if err != nil {
		return nil, err
	}

	replacementCount, err := performReplacement(mode, input, output, replMap)
	if err != nil {
		return nil, fmt.Errorf("replace operation failed: %w", err)
	}

	return map[string]interface{}{
		"path":              output,
		"replacements_made": replacementCount,
	}, nil
}

// validateAndExtractParams validates and extracts parameters from the params map.
func validateAndExtractParams(params map[string]interface{}) (input, output string, replMap map[string]string, mode string, err error) {
	input, ok := params["input"].(string)
	if !ok {
		return "", "", nil, "", fmt.Errorf("input parameter is required")
	}

	output, ok = params["output"].(string)
	if !ok {
		return "", "", nil, "", fmt.Errorf("output parameter is required")
	}

	replacements, ok := params["replacements"].(map[string]interface{})
	if !ok {
		return "", "", nil, "", fmt.Errorf("replacements parameter is required")
	}

	replMap = convertReplacements(replacements)

	mode, _ = params["mode"].(string)
	if mode == "" {
		mode = "id"
	}

	return input, output, replMap, mode, nil
}

// convertReplacements converts a map[string]interface{} to map[string]string.
func convertReplacements(replacements map[string]interface{}) map[string]string {
	replMap := make(map[string]string)
	for k, v := range replacements {
		replMap[k] = fmt.Sprintf("%v", v)
	}
	return replMap
}

// performReplacement performs the actual replacement based on the mode.
func performReplacement(mode, input, output string, replacements map[string]string) (int, error) {
	switch mode {
	case "id":
		return replaceByID(input, output, replacements)
	case "placeholder":
		return replacePlaceholders(input, output, replacements)
	default:
		return 0, fmt.Errorf("unsupported mode: %s (supported: id, placeholder)", mode)
	}
}

// replaceByID replaces text content in XML elements by ID attribute.
// This works for SVG and any XML file with id attributes on text elements.
func replaceByID(inputPath, outputPath string, replacements map[string]string) (int, error) {
	content, err := readXMLFile(inputPath)
	if err != nil {
		return 0, err
	}

	tokens, count, err := parseAndReplaceTokens(content, replacements)
	if err != nil {
		return 0, err
	}

	if err := writeXMLFile(outputPath, tokens); err != nil {
		return 0, err
	}

	return count, nil
}

// readXMLFile reads an XML file and returns its content.
func readXMLFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}
	return content, nil
}

// parseAndReplaceTokens parses XML content and replaces text in elements with matching IDs.
func parseAndReplaceTokens(content []byte, replacements map[string]string) ([]xml.Token, int, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(content)))
	var tokens []xml.Token
	var elementStack []xml.Token
	var targetID string
	replacementCount := 0

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse XML: %w", err)
		}

		token = xml.CopyToken(token)
		processedToken, replaced := processToken(token, &elementStack, &targetID, replacements)
		if replaced {
			replacementCount++
		}
		tokens = append(tokens, processedToken)
	}

	return tokens, replacementCount, nil
}

// processToken processes a single XML token, updating state and performing replacements.
func processToken(token xml.Token, elementStack *[]xml.Token, targetID *string, replacements map[string]string) (xml.Token, bool) {
	switch t := token.(type) {
	case xml.StartElement:
		*elementStack = append(*elementStack, token)
		*targetID = findTargetID(*elementStack, replacements)
		return token, false

	case xml.EndElement:
		if len(*elementStack) > 0 {
			*targetID = clearTargetIfExiting(*elementStack, *targetID)
			*elementStack = (*elementStack)[:len(*elementStack)-1]
		}
		return token, false

	case xml.CharData:
		return replaceCharData(t, *targetID, replacements)

	default:
		return token, false
	}
}

// findTargetID finds the nearest element ID in the stack that needs replacement.
func findTargetID(elementStack []xml.Token, replacements map[string]string) string {
	for i := len(elementStack) - 1; i >= 0; i-- {
		if elem, ok := elementStack[i].(xml.StartElement); ok {
			for _, attr := range elem.Attr {
				if attr.Name.Local == "id" {
					if _, exists := replacements[attr.Value]; exists {
						return attr.Value
					}
				}
			}
		}
	}
	return ""
}

// clearTargetIfExiting clears the target ID if we're exiting the element that set it.
func clearTargetIfExiting(elementStack []xml.Token, currentTargetID string) string {
	if len(elementStack) == 0 || currentTargetID == "" {
		return currentTargetID
	}

	if elem, ok := elementStack[len(elementStack)-1].(xml.StartElement); ok {
		for _, attr := range elem.Attr {
			if attr.Name.Local == "id" && attr.Value == currentTargetID {
				return ""
			}
		}
	}
	return currentTargetID
}

// replaceCharData replaces character data if a target ID is active.
func replaceCharData(charData xml.CharData, targetID string, replacements map[string]string) (xml.Token, bool) {
	if targetID == "" {
		return charData, false
	}

	replacement, exists := replacements[targetID]
	if !exists {
		return charData, false
	}

	text := strings.TrimSpace(string(charData))
	if text == "" || text == "\n" {
		return charData, false
	}

	return xml.CharData(replacement + "\n"), true
}

// writeXMLFile writes XML tokens to a file with proper formatting.
func writeXMLFile(path string, tokens []xml.Token) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	encoder := xml.NewEncoder(outputFile)
	encoder.Indent("", "  ")

	for _, token := range tokens {
		if err := encoder.EncodeToken(token); err != nil {
			return fmt.Errorf("failed to encode token: %w", err)
		}
	}

	if err := encoder.Flush(); err != nil {
		return fmt.Errorf("failed to flush encoder: %w", err)
	}

	return nil
}

// replacePlaceholders replaces simple text placeholders in any text file.
// Placeholders can be any string, e.g., {{name}}, $NAME, etc.
func replacePlaceholders(inputPath, outputPath string, replacements map[string]string) (int, error) {
	// Read input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read input file: %w", err)
	}

	// Perform replacements
	result := string(content)
	replacementCount := 0

	for placeholder, value := range replacements {
		count := strings.Count(result, placeholder)
		result = strings.ReplaceAll(result, placeholder, value)
		replacementCount += count
	}

	// Write output file
	if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
		return 0, fmt.Errorf("failed to write output file: %w", err)
	}

	return replacementCount, nil
}
