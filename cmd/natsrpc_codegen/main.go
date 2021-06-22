package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var (
	src        = flag.String("s", "", "src file(the path must be relative to go.mod)")
	clientMode = flag.Int("cm", 0, "client mode(0:sync, 1:async, 2:both)")
	//dest       = flag.String("d", "", "dest file")
	//outPackage = flag.String("op", "", "out package name")
)

// TmplParam 参数模板
type TmplParam struct {
	Name string
	Type string
}

// TmplMethod 方法模板
type TmplMethod struct {
	Name    string
	Comment string
	Param   []TmplParam
}

// TmplService 服务模板
type TmplService struct {
	Name    string
	Comment string
	Method  []TmplMethod
}

// Tmpl 模板
type Tmpl struct {
	OutPackage string
	Package    string
	Imports    []string
	Service    []TmplService
	ClientMode int
}

// GenText 模板执行
func GenText(tmpText string, data interface{}) ([]byte, error) {
	classTpl, err := template.New("temp").Parse(tmpText)
	if nil != err {
		return []byte(""), err
	}
	w := &bytes.Buffer{}
	if err := classTpl.Execute(w, data); nil != err {
		return []byte(""), err
	}
	return w.Bytes(), nil
}

// GenFile 生产代码文件
func GenFile(tmpText string, data interface{}, file string) error {
	src, err := GenText(tmpText, data)
	if nil != err {
		return err
	}

	if b, err := format.Source(src); nil != err {
		if errFile := ioutil.WriteFile(file, src, os.ModePerm); errFile != nil {
			return errFile
		}
		return err
	} else {
		return ioutil.WriteFile(file, b, os.ModePerm)
	}
}

func main() {
	flag.Parse()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *src, nil, 0)
	if err != nil {
		panic(err)
	}

	tmpl := Tmpl{
		OutPackage: f.Name.Name,
		Package:    f.Name.Name,
		ClientMode: *clientMode,
	}

	//tmpl.Imports = append(tmpl.Imports, fmt.Sprintf("\"%s/%s\"", *inPackage, (*src)[:strings.LastIndex(*src, "/")]))
	for _, v := range f.Imports {
		tmpl.Imports = append(tmpl.Imports, v.Path.Value)
	}

	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		if g.Tok != token.TYPE {
			continue
		}
		s := g.Specs[0].(*ast.TypeSpec)
		i := s.Type.(*ast.InterfaceType)

		ts := TmplService{
			Name:    s.Name.Name,
			Comment: s.Comment.Text(),
		}
		for _, m := range i.Methods.List {
			tm := TmplMethod{
				Name:    m.Names[0].Name,
				Comment: m.Comment.Text(),
			}
			f := m.Type.(*ast.FuncType)
			for i, p := range f.Params.List {
				tp := TmplParam{}
				if se, ok := p.Type.(*ast.SelectorExpr); ok {
					tp.Name = p.Names[0].Name
					tp.Type = se.X.(*ast.Ident).Name + "." + se.Sel.Name
				} else if se, ok := p.Type.(*ast.StarExpr); ok {
					tp.Name = p.Names[0].Name
					x := se.X.(*ast.SelectorExpr)
					tp.Type = x.X.(*ast.Ident).Name + "." + x.Sel.Name
				} else if se, ok := p.Type.(*ast.FuncType); ok {
					tp.Name = "reply"
					x := se.Params.List[0].Type.(*ast.StarExpr).X.(*ast.SelectorExpr)
					tp.Type = x.X.(*ast.Ident).Name + "." + x.Sel.Name
				}
				if i == 0 {
					if tp.Type != "context.Context" {
						panic("first param must be context.Context")
					}
				}
				tm.Param = append(tm.Param, tp)
			}
			if f.Results != nil {
				for _, p := range f.Results.List {
					tp := TmplParam{}
					if se, ok := p.Type.(*ast.SelectorExpr); ok {
						tp.Name = "rep"
						tp.Type = se.X.(*ast.Ident).Name + "." + se.Sel.Name
					} else if se, ok := p.Type.(*ast.StarExpr); ok {
						tp.Name = "rep"
						x := se.X.(*ast.SelectorExpr)
						tp.Type = x.X.(*ast.Ident).Name + "." + x.Sel.Name
					}
					tm.Param = append(tm.Param, tp)
					break
				}
			}

			ts.Method = append(ts.Method, tm)
		}
		tmpl.Service = append(tmpl.Service, ts)
	}

	dest := strings.Replace(*src, ".go", "", -1) + ".natsrpc.go"
	if err := GenFile(tFile, tmpl, dest); nil != err {
		panic(err)
	}

}
