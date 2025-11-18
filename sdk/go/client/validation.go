package client

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// Validator валидирует запросы и ответы по JSON Schema
type Validator struct {
	schemas map[string]*gojsonschema.Schema
}

// NewValidator создает новый валидатор
func NewValidator() *Validator {
	return &Validator{
		schemas: make(map[string]*gojsonschema.Schema),
	}
}

// LoadSchema загружает JSON Schema из файла
func (v *Validator) LoadSchema(name, schemaPath string) error {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	v.schemas[name] = schema
	return nil
}

// ValidateRequest валидирует запрос по схеме
func (v *Validator) ValidateRequest(schemaName string, data interface{}) error {
	schema, ok := v.schemas[schemaName]
	if !ok {
		// Если схема не найдена, пропускаем валидацию
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	documentLoader := gojsonschema.NewBytesLoader(jsonData)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		errors := ""
		for _, desc := range result.Errors() {
			errors += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("validation failed:\n%s", errors)
	}

	return nil
}

// ValidateResponse валидирует ответ по схеме
func (v *Validator) ValidateResponse(schemaName string, body io.Reader) error {
	schema, ok := v.schemas[schemaName]
	if !ok {
		// Если схема не найдена, пропускаем валидацию
		return nil
	}

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	documentLoader := gojsonschema.NewBytesLoader(bodyBytes)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		errors := ""
		for _, desc := range result.Errors() {
			errors += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("validation failed:\n%s", errors)
	}

	// Восстанавливаем body для дальнейшей обработки
	return nil
}

