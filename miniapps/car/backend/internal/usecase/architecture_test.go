package usecase_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

const modulePath = "github.com/enoro/eiger-webapp/miniapps/car/backend"

// usecase 層は、repository の実装・handler・infrastructure・テスト用部品を import しない。
// 依存してよいのは、標準ライブラリ・domain・usecase 配下だけ(repository は domain/repository の interface のみ)。
func TestUsecase層は外側の層をimportしない(t *testing.T) {
	forbidden := []string{
		modulePath + "/internal/repository",
		modulePath + "/internal/handler",
		modulePath + "/internal/infrastructure",
		modulePath + "/internal/testutil",
	}
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
			for _, bad := range forbidden {
				if p == bad || strings.HasPrefix(p, bad+"/") {
					t.Errorf("%s が外側の層を import している: %s", path, p)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
