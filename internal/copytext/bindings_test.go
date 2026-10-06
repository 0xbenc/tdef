package copytext

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// Verify even screens/branches not exercised by rendering tests bind real copy
// keys and exactly the placeholders the authored template expects.
func TestProductionBindings(t *testing.T) {
	used := map[string]bool{}
	for _, area := range []string{"game", "render", "tui"} {
		dir := filepath.Join("..", "..", area)
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
				continue
			}
			fs := token.NewFileSet()
			file, err := parser.ParseFile(fs, filepath.Join(dir, entry.Name()), nil, 0)
			if err != nil {
				t.Fatal(err)
			}
			ast.Inspect(file, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				selector, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				pkg, ok := selector.X.(*ast.Ident)
				if !ok || pkg.Name != "copytext" {
					return true
				}
				literal := func(expr ast.Expr) (string, bool) {
					str, ok := expr.(*ast.BasicLit)
					if !ok || str.Kind != token.STRING {
						return "", false
					}
					value, err := strconv.Unquote(str.Value)
					return value, err == nil
				}
				key, ok := literal(call.Args[0])
				if !ok {
					// The overworld ledger chooses a display label by known floor ID.
					return true
				}
				used[key] = true
				want, exists := required[key]
				if !exists {
					t.Errorf("%s: unknown copy key %q", fs.Position(call.Pos()), key)
					return true
				}
				var names []string
				if selector.Sel.Name == "Format" {
					if len(call.Args)%2 != 1 {
						t.Errorf("%s: incomplete placeholder pairs", key)
						return true
					}
					for i := 1; i < len(call.Args); i += 2 {
						name, ok := literal(call.Args[i])
						if !ok {
							t.Errorf("%s: placeholder names must be explicit", key)
						}
						names = append(names, name)
					}
				}
				slices.Sort(names)
				if !slices.Equal(names, want) {
					t.Errorf("%s: binds %v, template needs %v", key, names, want)
				}
				return true
			})
		}
	}
	for _, floor := range []string{"rift", "halls", "garden", "rotunda"} {
		used["places."+floor+".ledger_name"] = true
	}
	for key := range required {
		if !used[key] {
			t.Errorf("%s: copy entry has no production binding", key)
		}
	}
}
