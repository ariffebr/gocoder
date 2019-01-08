package gocoder

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type GoStruct struct {
	rootExpr *GoExpr

	astExpr *ast.StructType

	goFields []*GoField

	methods []*GoFunc

	spec *ast.TypeSpec

	structTag string
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

func (p *GoStruct) Tag() StructTag {
	decls := p.rootExpr.astFile.Decls

	for _, decl := range decls {

		if dcl, ok := decl.(*ast.GenDecl); ok {
			specs := dcl.Specs
			for _, spec := range specs {

				if s, ok2 := spec.(*ast.TypeSpec); ok2 {
					if s.Name.Name == p.Name() {
						fmt.Println("Stuct Found!!! ---------> ", s.Name)

						docList := dcl.Doc.List

						var tags []string
						for _, doc := range docList {

							if strings.Contains(doc.Text, "`") && strings.Count(doc.Text, "`") > 1 {
								tag := doc.Text[strings.Index(doc.Text, "`") : strings.LastIndex(doc.Text, "`")+1]
								tag = strings.Trim(tag, "\"")
								tag = strings.Trim(tag, "`")

								for {
									if !strings.HasSuffix(tag, ";") {
										break
									}
									fmt.Println("has suffix semicolon")
									tag = tag[0 : len(tag)-1]
								}
								tag = strings.TrimSpace(tag)
								tags = append(tags, tag)
							}
						}

						if len(tags) == 0 {
							break
						}
						return StructTag(strings.Join(tags, ";"))
					}
				}
			}

		}
	}

	return StructTag("")
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
