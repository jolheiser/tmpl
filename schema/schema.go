package schema

import (
	_ "embed"
	"fmt"
	"io"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var (
	//go:embed tmpl.json
	schema       []byte
	schemaLoader = gojsonschema.NewBytesLoader(schema)
)

// Lint is for linting a recipe against the schema
func Lint(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	var m any
	if err := Unmarshal(data, &m); err != nil {
		return err
	}
	sourceLoader := gojsonschema.NewGoLoader(m)
	result, err := gojsonschema.Validate(schemaLoader, sourceLoader)
	if err != nil {
		return err
	}
	if len(result.Errors()) > 0 {
		return ResultErrors(result.Errors())
	}
	return nil
}

// ResultErrors is a slice of gojsonschema.ResultError that implements error
type ResultErrors []gojsonschema.ResultError

// Error implements error
func (r ResultErrors) Error() string {
	errs := make([]string, 0, len(r))
	for _, re := range r {
		errs = append(errs, fmt.Sprintf("%s: %s", re.Field(), re.Description()))
	}
	return strings.Join(errs, " | ")
}
