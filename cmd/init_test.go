package cmd

import (
	"strings"
	"testing"

	"go.jolheiser.com/tmpl/schema"

	"github.com/matryer/is"
)

func TestInitSchema(t *testing.T) {
	assert := is.New(t)

	err := schema.Lint(strings.NewReader(initConfig))
	assert.NoErr(err) // Init config should conform to schema
}
