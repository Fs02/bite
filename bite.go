package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	dt = data{
		Pkg:   "main",
		Var:   "files",
		Paths: make(map[string]string),
	}
	output   = ""
	trimPath = ""
)

func main() {
	flag.StringVar(&dt.Pkg, "package", dt.Pkg, "`name` of the package to generate")
	flag.StringVar(&dt.Var, "var", dt.Var, "`name` of the variable to generate")
	flag.StringVar(&output, "output", output, "`filename` to write the output to")
	flag.StringVar(&trimPath, "trim", trimPath, "path `prefix` to remove from the resulting file path")
	flag.Parse()

	if output == "" {
		flag.PrintDefaults()
		log.Fatal("-output is required.")
	}

	for _, g := range flag.Args() {
		fmt.Println(g)
		matches, err := filepath.Glob(g)
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range matches {
			dt.Paths[m] = strings.TrimPrefix(m, trimPath)
		}
	}

	file, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, dt)
	if err != nil {
		log.Fatal(err)
	}
}

type data struct {
	Pkg   string
	Var   string
	Paths map[string]string
}

func content(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%#v,\n", content), nil
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("file").Funcs(template.FuncMap{"content": content}).Parse(tmplString))
}

const tmplString = `// Package {{ .Pkg }} is generated by github.com/Fs02/bite
package {{ .Pkg }}

var {{ .Var }} = map[string][]byte{
	{{range $path, $name := .Paths }}"{{ $name }}": {{ content $path }}{{ end }}}
`
