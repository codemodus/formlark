package main

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Templates struct {
	m          map[string]*template.Template
	delimLeft  string
	delimRight string
}

func NewTemplates(delimLeft, delimRight string) *Templates {
	m := make(map[string]*template.Template)
	if delimLeft == "" {
		delimLeft = "{{"
	}
	if delimRight == "" {
		delimRight = "}}"
	}
	t := &Templates{m, delimLeft, delimRight}
	return t
}

func (t *Templates) ParseDir(dir string) {
	fMap := template.FuncMap{}

	filepath.Walk(dir, func(p string, i os.FileInfo, e error) error {
		if i == nil || i.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ".tmpl.html") || !strings.HasSuffix(p, ".html") {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}

		r, err := t.processTemplate(dir, rel)
		if err != nil {
			panic(err)
		}
		b := &bytes.Buffer{}
		if _, err = b.ReadFrom(r); err != nil {
			panic(err)
		}

		n := filepath.ToSlash(rel)
		t.m[n] = template.New(``).Funcs(fMap)
		template.Must(t.m[n].Parse(b.String()))

		if _, err = os.Stat(path.Join(dir, rel[:len(rel)-5], "index.html")); os.IsNotExist(err) {
			t.m[n[:len(n)-5]] = t.m[n]
		}

		if rel[len(rel)-10:] == "index.html" {
			t.m[strings.TrimRight(n[:len(n)-10], "/")] = t.m[n]
		}
		return nil
	})
}

func (t *Templates) processTemplate(dir, relPath string) (io.Reader, error) {
	iCt := len(t.delimLeft) + len("include")

	f, err := os.Open(path.Join(dir, relPath))
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(f)
	b := &bytes.Buffer{}

	for s.Scan() {
		cs := strings.Replace(s.Text(), " ", "", -1)
		if len(cs) > iCt && cs[:iCt] == t.delimLeft+"include" {
			if len(cs) < iCt+1+2 {
				return nil, errors.New("malformed include")
			}
			part := cs[iCt+1:]
			i := strings.Index(part, `"`)
			if i < 0 {
				return nil, errors.New("malformed include")
			}

			np := part[:i]
			if np[0:1] != "/" {
				np = path.Join(path.Dir(relPath), np)
			}

			r, err := t.processTemplate(dir, np)
			if err != nil {
				return nil, err
			}
			if _, err = b.ReadFrom(r); err != nil {
				return nil, err
			}

			b.Write([]byte("\n"))
			continue
		}

		b.Write(s.Bytes())
		b.Write([]byte("\n"))
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return b, nil
}

func (t *Templates) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	l, r := 0, 0
	if name[0:1] == `/` {
		l = 1
	}
	if name[len(name)-1:] == `/` {
		r = 1
	}
	name = name[l : len(name)-r]

	err := t.m[name].ExecuteTemplate(wr, ``, data)
	return err
}
