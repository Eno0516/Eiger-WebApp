package handler_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

const modulePath = "github.com/enoro/eiger-webapp/miniapps/car/backend"

// handler 層が呼ぶのは usecase だけ。repository(interface・実装とも)や infrastructure、テスト用部品には依存しない。
func TestHandler層はusecase以外の内側と外側をimportしない(t *testing.T) {
	forbidden := []string{
		modulePath + "/internal/domain/repository",
		modulePath + "/internal/repository",
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
					t.Errorf("%s が handler から呼んではいけない層を import している: %s", path, p)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
