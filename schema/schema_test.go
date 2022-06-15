package schema

import (
	"embed"
	"errors"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

//go:embed testdata
var testdata embed.FS

func TestSchema(t *testing.T) {
	tt := []struct {
		Name   string
		NumErr int
	}{
		{Name: "good", NumErr: 0},
		{Name: "bad", NumErr: 10},
		{Name: "empty", NumErr: 1},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assert := is.New(t)

			fi, err := testdata.Open(fmt.Sprintf("testdata/%s.yaml", tc.Name))
			assert.NoErr(err) // Should open test file

			err = Lint(fi)
			if tc.NumErr > 0 {
				var rerrs ResultErrors
				assert.True(errors.As(err, &rerrs))  // Error should be ResultErrors
				assert.True(len(rerrs) == tc.NumErr) // Number of errors should match test case
			} else {
				assert.NoErr(err) // Good schemas shouldn't return errors
			}
		})
	}
}
