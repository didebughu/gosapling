package pointer

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `check nullpointer.`

var Analyzer = &analysis.Analyzer{
	Name:     "checkpointer",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var name, paramorder *bool

func init() {
	name = Analyzer.Flags.Bool("name", true, "checkpointer")
	paramorder = Analyzer.Flags.Bool("paramorder", true, "checkpointer")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	pointers := []Pointer{}
	// ast.Print(pass.Fset, pass.Files)
	inspect.Preorder(nil, func(n ast.Node) {
		if genDecl, OK := n.(*ast.GenDecl); OK {
			for _, decl := range genDecl.Specs {
				if valueSpec, ok := decl.(*ast.ValueSpec); ok {
					if _, okk := valueSpec.Type.(*ast.StarExpr); okk && valueSpec.Values == nil {
						for _, name := range valueSpec.Names {
							pointers = append(pointers,
								Pointer{name: name.Name, namePos: name.NamePos, isNil: true})
						}
					}
				}
			}
		}
		if assignStmt, OK := n.(*ast.AssignStmt); OK {
			lname := getAssignStmtLhsNames(assignStmt)
			for _, name := range lname {
				for i := range pointers {
					if pointers[i].name == name {
						pointers[i].isNil = false
					}
				}
			}
		}
		if exprStmt, OK := n.(*ast.ExprStmt); OK {
			var tmp []Pointer
			ast.Inspect(exprStmt, func(node ast.Node) bool {
				starExpr, ok := node.(*ast.StarExpr)
				if ok {
					tmp = append(tmp,
						Pointer{name: starExpr.X.(*ast.Ident).Name, namePos: starExpr.X.(*ast.Ident).Pos(), isNil: true})
					return false
				}
				return true
			})
			for i := range tmp {
				for j := range pointers {
					if pointers[j].name == tmp[i].name && pointers[j].isNil {
						fmt.Println(pass.Fset.Position(tmp[i].namePos), "(warning) Possible null pointer dereference: ", pointers[j].name)
					}
				}
			}

		}
	})

	return nil, nil
}

func getAssignStmtLhsNames(assign *ast.AssignStmt) []string {
	var names []string
	for _, expr := range assign.Lhs {
		ident, ok := expr.(*ast.Ident)
		if ok {
			names = append(names, ident.Name)
		}
	}
	return names
}
