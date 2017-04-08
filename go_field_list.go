package gocoder

import (
	"context"
	"go/ast"
	"go/token"
)

type GoFieldList struct {
	rootExpr *GoExpr
	astExpr  *ast.FieldList

	goFields []*GoField
}

func newFieldList(rootExpr *GoExpr, astFieldList *ast.FieldList) *GoFieldList {
	gfl := &GoFieldList{
		rootExpr: rootExpr,
		astExpr:  astFieldList,
	}

	gfl.load()

	return gfl
}

func (p *GoFieldList) Print() error {
	return ast.Print(p.rootExpr.astFileSet, p.astExpr)
}

func (p *GoFieldList) NumFields() int {
	return p.astExpr.NumFields()
}

func (p *GoFieldList) Field(i int) *GoField {
	return p.goFields[i]
}

func (p *GoFieldList) TypesAre(paramsType ...string) bool {

	if p.astExpr.NumFields() != len(paramsType) {
		return false
	}

	for i := 0; i < p.astExpr.NumFields(); i++ {
		if !isFieldTypeOf(p.astExpr.List[i], paramsType[i]) {
			return false
		}
	}

	return true
}

func (p *GoFieldList) Position() token.Position {
	return p.rootExpr.astFileSet.Position(p.astExpr.Pos())
}

func (p *GoFieldList) load() {
	for i := 0; i < len(p.astExpr.List); i++ {
		p.goFields = append(p.goFields, newGoField(p.rootExpr, p.astExpr.List[i]))
	}
}

func (p *GoFieldList) Inspect(f InspectFunc, ctx context.Context) {
	for i := 0; i < len(p.goFields); i++ {
		p.goFields[i].Inspect(f, ctx)
	}
}

func (p *GoFieldList) goNode() {}
