package gocoder

import (
	"go/ast"
	"go/token"
)

type GoStruct struct {
	rootExpr *GoExpr

	astExpr *ast.StructType

	goFields []*GoField

	methods []*GoFunc

	spec *ast.TypeSpec
}

func newGoStruct(rootExpr *GoExpr, spec *ast.TypeSpec, expr *ast.StructType, options ...Option) *GoStruct {

	g := &GoStruct{
		rootExpr: rootExpr,
		astExpr:  expr,
		spec:     spec,
	}

	g.load()

	return g
}

func (p *GoStruct) Name() string {
	if p.spec == nil {
		return ""
	}

	return p.spec.Name.String()
}

func (p *GoStruct) NumFields() int {
	return p.astExpr.Fields.NumFields()
}

func (p *GoStruct) Field(i int) *GoField {
	return p.goFields[i]
}

func (p *GoStruct) NumMethod() int {
	return len(p.methods)
}

func (p *GoStruct) Method(i int) *GoFunc {
	return p.methods[i]
}

func (p *GoStruct) Position() (token.Position, token.Position) {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos()), p.rootExpr.astFileSet.Position(p.astExpr.End())
}

func (p *GoStruct) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoStruct) load() {

	var goFileds []*GoField

	for i := 0; i < len(p.astExpr.Fields.List); i++ {
		field := p.astExpr.Fields.List[i]
		goFileds = append(goFileds, newGoField(p.rootExpr, field))
	}

	p.goFields = goFileds

	p.loadFuncs()
}

func (p *GoStruct) loadFuncs() {

	goPckg := p.rootExpr.options.GoPackage

	if goPckg != nil {
		for i := 0; i < goPckg.NumFuncs(); i++ {
			fn := goPckg.Func(i)

			if fn.Receiver() == p.Name() {
				p.methods = append(p.methods, fn)
			}
		}

	}
}

func (p *GoStruct) goNode() {}
