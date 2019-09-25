package loaderx

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDocTextOfNode(node ast.Node) string {
	var doc *ast.CommentGroup

	switch node.(type) {
	case *ast.File:
		doc = node.(*ast.File).Doc
	case *ast.GenDecl:
		doc = node.(*ast.GenDecl).Doc
	case *ast.FuncDecl:
		doc = node.(*ast.FuncDecl).Doc
	case *ast.Field:
		doc = node.(*ast.Field).Doc
	case *ast.ImportSpec:
		doc = node.(*ast.ImportSpec).Doc
	case *ast.TypeSpec:
		doc = node.(*ast.TypeSpec).Doc
	case *ast.ValueSpec:
		doc = node.(*ast.ValueSpec).Doc
	}

	if doc != nil {
		return strings.TrimSpace(doc.Text())
	}
	return ""
}

func TestCommentScanner(t *testing.T) {
	tt := assert.New(t)

	fileContent, _ := ioutil.ReadFile("./fixtures/comments.go")
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "./fixtures/comments.go", string(fileContent), parser.ParseComments)
	commentScanner := NewCommentScanner(fset, file)

	ast.Inspect(file, func(node ast.Node) bool {
		comments := commentScanner.CommentsOf(node)

		switch node.(type) {
		case *ast.File, *ast.Field, ast.Decl:
			tt.Equal(getDocTextOfNode(node), comments)
		case *ast.ImportSpec:
			c := getDocTextOfNode(node)
			if c == "" {
				c = "import"
			}
			tt.Equal(c, comments)
		case *ast.TypeSpec:
			c := getDocTextOfNode(node)
			if c == "" {
				c = "type"
			}
			tt.Equal(c, comments)
		case *ast.ValueSpec:
			c := getDocTextOfNode(node)
			if c == "" {
				c = "var"
			}
			tt.Equal(c, comments)
		case ast.Stmt:
			commentGroupList := commentScanner.CommentMap[node]
			tt.Equal(StringifyCommentGroup(commentGroupList...), comments)
		}
		return true
	})
}
