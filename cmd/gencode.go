package main

import (
	"flag"
	"fmt"
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
	inPackage = flag.String("ip", "", "in package name")
	src         = flag.String("s", "", "src file(the path must be relative to go.mod)")
	dest        = flag.String("d", "", "dest file")
	outPackage = flag.String("op", "", "out package name")
)

type TmplParam struct {
	Name string
	Type  string
}

type TmplMethod struct {
	Name  string
	Comment string
	Param []TmplParam
}

type TmplService struct {
	Name   string
	Comment string
	Method []TmplMethod
}

type Tmpl struct {
	OutPackage string
	Package string
	Imports []string
	Service []TmplService
}

// ExecTemplate 模板执行
func GenText(tmpText string, data interface{}) (string, error) {
	classTpl, err := template.New("temp").Parse(tmpText)
	if nil != err {
		return "", err
	}
	w := &strings.Builder{}
	if err := classTpl.Execute(w, data); nil != err {
		return "", err
	}
	return w.String(), nil
}
// GenFile 生产代码文件
func GenFile(tmpText string, data interface{}, file string) error {
	src, err := GenText(tmpText, data)
	if nil != err {
		return err
	}

	b, err := format.Source([]byte(src))
	if nil != err {
		fmt.Println(err)
		//return err
	}

	if err := ioutil.WriteFile(file, b, os.ModePerm); nil != err {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, *src, nil, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	//ast.Print(fset, f)

	tmpl := Tmpl{
		OutPackage: *outPackage,
		Package:f.Name.Name,
	}

	tmpl.Imports = append(tmpl.Imports,fmt.Sprintf("\"%s/%s\"",*inPackage,(*src)[:strings.LastIndex(*src,"/")]))
	for _,v := range f.Imports {
		tmpl.Imports = append(tmpl.Imports,v.Path.Value)
	}

	for _, d := range f.Decls {
		g,ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		if g.Tok !=token.TYPE {
			continue
		}
		s :=g.Specs[0].(*ast.TypeSpec)
		i := s.Type.(*ast.InterfaceType)

		ts := TmplService{
			Name:s.Name.Name,
			Comment:s.Comment.Text(),
		}
		for _,m := range i.Methods.List {
			tm := TmplMethod{
				Name: m.Names[0].Name,
				Comment:m.Comment.Text(),
			}
			for i,p:= range m.Type.(*ast.FuncType).Params.List {
				tp := TmplParam{}
				if se,ok := p.Type.(*ast.SelectorExpr); ok {
					tp.Name = p.Names[0].Name
					tp.Type= se.X.(*ast.Ident).Name+"."+se.Sel.Name
				} else if se,ok := p.Type.(*ast.StarExpr); ok  {
					tp.Name = p.Names[0].Name
					x := se.X.(*ast.SelectorExpr)
					tp.Type= x.X.(*ast.Ident).Name+"."+x.Sel.Name
				}
				if i==0 {
					if tp.Type!="context.Context" {
						panic("first param must be context.Context")
					}
				}
				tm.Param = append(tm.Param, tp)
			}
			ts.Method = append(ts.Method, tm)
		}
		tmpl.Service = append(tmpl.Service,ts)
	}

	if err:=GenFile(tFile,tmpl,*dest); nil!=err {
		panic(err)
	}

}
