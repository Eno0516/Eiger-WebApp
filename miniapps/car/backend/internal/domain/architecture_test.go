package domain_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

const modulePath = "github.com/enoro/eiger-webapp/miniapps/car/backend"

// domain 層は、標準ライブラリと domain 配下のパッケージ以外を import しない。
// usecase・handler・repository 実装・本体(backend)など、外側の層に依存すると失敗する。
func TestDomain層は外側の層をimportしない(t *testing.T) {
	fset := token.NewFileSet()
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			if isStdlib(p) || strings.HasPrefix(p, modulePath+"/internal/domain") {
				continue
			}
			t.Errorf("%s が domain の外を import している: %s", path, p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// 標準ライブラリは、import パスの先頭要素にドットを含まない。
func isStdlib(path string) bool {
	first, _, _ := strings.Cut(path, "/")
	return !strings.Contains(first, ".")
}
