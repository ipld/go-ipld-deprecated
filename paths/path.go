package paths

import (
	"strings"
)

// Escape path component. The special characters ("@" and "\") are escaped to
// allow mixing the path component with directives (starting with "@") in IPLD
// data structure.
func EscapePathComponent(comp string) string {
	comp = strings.Replace(comp, "\\", "\\\\", -1)
	comp = strings.Replace(comp, "@", "\\@", -1)
	return comp
}

// Unescape path component from the IPLD data structure. Special characters are
// unescaped. See also EscapePathComponent function.
func UnescapePathComponent(comp string) string {
	comp = strings.Replace(comp, "\\@", "@", -1)
	comp = strings.Replace(comp, "\\\\", "\\", -1)
	return comp
}
