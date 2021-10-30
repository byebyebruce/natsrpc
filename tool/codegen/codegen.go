package codegen

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"
)

var tpl *template.Template

func init() {
	tmp, err := template.New("codegen").Parse(tempFile)
	if err != nil {
		panic(err)
	}
	tpl = tmp
}

// Template 模板
func Template() *template.Template {
	return tpl
}

// GenFile 生成代码文件
func GenFile(data FileSpec, file string) error {
	src, err := GenText(tpl, data)
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

// GenText 模板执行
func GenText(t *template.Template, data interface{}) ([]byte, error) {
	w := &bytes.Buffer{}
	if err := t.Execute(w, data); nil != err {
		return []byte(""), err
	}
	return w.Bytes(), nil
}
