package registry

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/huandu/xstrings"
)

var funcMap = map[string]any{
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
	"trim_prefix": func(in, trim string) string {
		return strings.TrimPrefix(in, trim)
	},
	"trim_suffix": func(in, trim string) string {
		return strings.TrimSuffix(in, trim)
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
