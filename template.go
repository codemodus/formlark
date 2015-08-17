package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Page struct {
	AppName  string
	URLLogin string
	Misc     string
}

func NewPage() *Page {
	return &Page{AppName: "Formlark", URLLogin: "/admin/login"}
}

func getTemplates() *template.Template {
	dir := `templates`
	fMap := template.FuncMap{}
	ts := template.New(``).Funcs(fMap)
	filepath.Walk(dir, func(p string, i os.FileInfo, e error) error {
		if i == nil || i.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}

		f, err := ioutil.ReadFile(p)
		if err != nil {
			panic(err)
		}

		template.Must(ts.New(filepath.ToSlash(rel)).Parse(string(f)))
		return nil
	})

	return ts
}
