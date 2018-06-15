package nursetags

import (
	"fmt"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"
)

func parseExpr(expr string) (Set, error) {
	tr, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	return eval(tr), nil
}

func eval(t ast.Expr) Set {
	switch t.(type) {
	case *ast.BinaryExpr:
		return evalBinaryExpr(t.(*ast.BinaryExpr))
	case *ast.ParenExpr:
		return evalParenExpr(t.(*ast.ParenExpr))
	case *ast.BasicLit:
		return evalBasicLit(t.(*ast.BasicLit))
	case *ast.Ident:
		return evalIdent(t.(*ast.Ident))
	}
	return nil
}

func evalBinaryExpr(t *ast.BinaryExpr) Set {
	switch t.Op {
	case token.AND:
		return eval(t.X).Intersect(eval(t.Y))
	case token.OR:
		return eval(t.X).Union(eval(t.Y))
	case token.XOR:
		return eval(t.X).Difference(eval(t.Y))
	}
	return nil
}

func evalParenExpr(t *ast.ParenExpr) Set {
	return eval(t.X)
}

func evalBasicLit(t *ast.BasicLit) Set {
	key := strings.Trim(t.Value, `"'`)
	if tag, ok := databaseTagRead(key); ok {
		return Set(tag)
	}
	return Set{}
}

func evalIdent(t *ast.Ident) Set {
	key := t.Name
	if tag, ok := databaseTagRead(key); ok {
		return Set(tag)
	}
	return Set{}
}

func repr(t ast.Expr) string {
	switch t.(type) {
	case *ast.BinaryExpr:
		return reprBinaryExpr(t.(*ast.BinaryExpr))
	case *ast.ParenExpr:
		return reprParenExpr(t.(*ast.ParenExpr))
	case *ast.BasicLit:
		return reprBasicLit(t.(*ast.BasicLit))
	case *ast.Ident:
		return reprIdent(t.(*ast.Ident))
	default:
		return fmt.Sprint(t)
	}
}

func reprBinaryExpr(t *ast.BinaryExpr) string {
	return fmt.Sprintf("%s op:%s %s", repr(t.X), t.Op, repr(t.Y))
}

func reprParenExpr(t *ast.ParenExpr) string {
	return fmt.Sprintf("paren:( %s )", repr(t.X))
}

func reprBasicLit(t *ast.BasicLit) string {
	return fmt.Sprintf("rune:%s", t.Value)
}

func reprIdent(t *ast.Ident) string {
	return fmt.Sprintf("ident:%s", t.Name)
}
