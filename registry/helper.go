package registry

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/huandu/xstrings"
)

var funcMap = map[string]interface{}{

	// String conversions
	"upper":  strings.ToUpper,
	"lower":  strings.ToLower,
	"title":  strings.Title,
	"snake":  xstrings.ToSnakeCase,
	"kebab":  xstrings.ToKebabCase,
	"pascal": xstrings.ToCamelCase,
	"camel": func(in string) string {
		return xstrings.FirstRuneToLower(xstrings.ToCamelCase(in))
	},

	// Other
	"env": os.Getenv,
	"sep": func() string {
		return string(filepath.Separator)
	},
	"time": func(fmt string) string {
		return time.Now().Format(fmt)
	},
}
