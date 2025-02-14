package handlers

import (
	"html/template"
	"os"
	"path/filepath"
)

var Templates *template.Template

func ParseTemplates() (*template.Template, error) {
    tmpl := template.New("")

    err := filepath.Walk("view", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && filepath.Ext(path) == ".html" {
            data, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            _, err = tmpl.New(path).Parse(string(data))
            if err != nil {
                return err
            }
        }
        return nil
    })

    return tmpl, err
}
