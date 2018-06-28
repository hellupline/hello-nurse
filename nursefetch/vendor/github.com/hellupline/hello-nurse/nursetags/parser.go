package nursetags

import (
	"fmt"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"
)

func (d *Database) ParseExpr(expr string) (Set, error) {
	tr, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	return d.eval(tr), nil
}

func (d *Database) eval(t ast.Expr) Set {
	switch t.(type) {
	case *ast.BinaryExpr:
		return d.evalBinaryExpr(t.(*ast.BinaryExpr))
	case *ast.ParenExpr:
		return d.evalParenExpr(t.(*ast.ParenExpr))
	case *ast.BasicLit:
		return d.evalBasicLit(t.(*ast.BasicLit))
	case *ast.Ident:
		return d.evalIdent(t.(*ast.Ident))
	}
	return nil
}

func (d *Database) evalBinaryExpr(t *ast.BinaryExpr) Set {
	switch t.Op {
	case token.AND:
		return d.eval(t.X).Intersect(d.eval(t.Y))
	case token.OR:
		return d.eval(t.X).Union(d.eval(t.Y))
	case token.XOR:
		return d.eval(t.X).Difference(d.eval(t.Y))
	}
	return nil
}

func (d *Database) evalParenExpr(t *ast.ParenExpr) Set {
	return d.eval(t.X)
}

func (d *Database) evalBasicLit(t *ast.BasicLit) Set {
	key := strings.Trim(t.Value, `"'`)
	if tag, ok := DefaultDatabase.TagRead(key); ok {
		return Set(tag)
	}
	return Set{}
}

func (d *Database) evalIdent(t *ast.Ident) Set {
	key := t.Name
	if tag, ok := DefaultDatabase.TagRead(key); ok {
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