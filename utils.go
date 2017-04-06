package gocoder

import (
	"go/ast"
)

func fieldTypeToStringType(field *ast.Field) string {
	exprStr := ""

	selector := 0
	mapChar := ""

	ast.Inspect(field.Type, func(n ast.Node) bool {
		switch ex := n.(type) {
		case *ast.SelectorExpr:
			{
				selector += 1
			}
		case *ast.Ident:

			exprStr += ex.Name
			if selector > 0 {
				exprStr += "."
				selector--
			}

			if len(mapChar) > 0 {
				exprStr += mapChar
				mapChar = ""
			}
		case *ast.InterfaceType:
			exprStr += "interface{}"
		case *ast.MapType:
			exprStr += "map["
			mapChar = "]"
		case *ast.StarExpr:
			exprStr += "*"
		case *ast.ArrayType:
			exprStr += "[]"
		case *ast.Ellipsis:
			exprStr += "..."
		}

		return true
	})

	return exprStr
}

func isFieldTypeOf(field *ast.Field, strType string) bool {
	return fieldTypeToStringType(field) == strType
}

func isCallingFuncOf(expr interface{}, name string) bool {
	switch ex := (expr).(type) {
	case *ast.CallExpr:
		{
			switch fn := ex.Fun.(type) {
			case *ast.SelectorExpr:
				{
					return fn.Sel.Name == name
				}
			case *ast.Ident:
				{
					return fn.Name == name
				}
			}
		}
	}
	return false
}

func trimStarExpr(expr ast.Expr) ast.Expr {

	typeExpr := expr

	for {
		starExpr, isStar := typeExpr.(*ast.StarExpr)

		if isStar {
			typeExpr = starExpr.X
			continue
		}

		break
	}

	return typeExpr
}
