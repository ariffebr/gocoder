package gocoder

import (
	"go/ast"
)

type GoSelector struct {
	*GoExpr
	rootExpr *GoExpr

	goSelIdent *GoIdent
	goXExpr    *GoExpr

	astExpr *ast.SelectorExpr
}

func newGoSelector(rootExpr *GoExpr, astSelector *ast.SelectorExpr) *GoSelector {
	g := &GoSelector{
		rootExpr: rootExpr,
		astExpr:  astSelector,
		GoExpr:   newGoExpr(rootExpr, astSelector),
	}

	g.load()

	return g
}

func (p *GoSelector) load() {
	if p.astExpr.X != nil {
		p.goXExpr = newGoExpr(p.rootExpr, p.astExpr.X)
	}

	if p.astExpr.Sel != nil {
		p.goSelIdent = newGoIdent(p.rootExpr, p.astExpr.Sel)
	}
}

func (p *GoSelector) Inspect(f func(GoNode) bool) {
	p.goXExpr.Inspect(f)
	p.goSelIdent.Inspect(f)
}

func (p *GoSelector) goNode() {}
